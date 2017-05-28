package network

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/quilt/quilt/db"
	"github.com/quilt/quilt/join"
	"github.com/quilt/quilt/stitch"
	"github.com/quilt/quilt/util"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-iptables/iptables"
	"github.com/vishvananda/netlink"
)

// IPTables is an interface to *iptables.IPTables.
type IPTables interface {
	Append(string, string, ...string) error
	AppendUnique(string, string, ...string) error
	Delete(string, string, ...string) error
	List(string, string) ([]string, error)
}

func runNat(conn db.Conn) {
	tables := []db.TableType{db.ContainerTable, db.ConnectionTable, db.MinionTable}
	for range conn.TriggerTick(30, tables...).C {
		minion := conn.MinionSelf()
		if minion.Role != db.Worker {
			continue
		}

		connections := conn.SelectFromConnection(nil)
		containers := conn.SelectFromContainer(func(c db.Container) bool {
			return c.IP != ""
		})

		ipt, err := iptables.New()
		if err != nil {
			log.WithError(err).Error("Failed to get iptables handle")
			continue
		}

		if err := updateNAT(ipt, containers, connections); err != nil {
			log.WithError(err).Error("Failed to update NAT rules")
		}
	}
}

// updateNAT sets up iptables rules of three categories:
// "default rules" are general rules that must be in place for the PREROUTING
// rules to work. When syncing "default rules" we don't remove any other rules
// that may be in place.
// "prerouting rules" are responsible for routing traffic to specific
// containers. They overwrite any pre-existing or outdated rules.
// "postrouting rules" are responsible for routing traffic from containers
// to the public internet. They overwrite any pre-existing or outdated rules.
func updateNAT(ipt IPTables, containers []db.Container,
	connections []db.Connection) error {

	publicInterface, err := getPublicInterface()
	if err != nil {
		return fmt.Errorf("get public interface: %s", err)
	}

	if err := setDefaultRules(ipt); err != nil {
		return err
	}

	prerouting := preroutingRules(publicInterface, containers, connections)
	if err := syncChain(ipt, "nat", "PREROUTING", prerouting); err != nil {
		return err
	}

	postrouting := postroutingRules(publicInterface, containers, connections)
	return syncChain(ipt, "nat", "POSTROUTING", postrouting)
}

var flagRegex = regexp.MustCompile(`-{1,2}(\S+) (\S+)(.*)`)

// ruleKey transforms an iptables rule into a string that is consistent between
// changes to the order of options.
// It handles rules of the form `-k v --another val -j target -with flags`. In
// other words, we parse out key value pairs as denoted by hyphens, unless
// we've reached the jump flag, in which case the value is the remaining
// string.
// The algorithm works by matching three parts: the flag name, the flag
// value, and the remaining string. The remaining string will begin with
// the next flag, and is processed on the next iteration.
func ruleKey(intf interface{}) interface{} {
	opts := map[string]string{}
	remaining := intf.(string)
	for remaining != "" {
		matches := flagRegex.FindStringSubmatch(remaining)
		if matches == nil {
			log.WithField("rule", intf.(string)).Error("Failed to parse rule")
			return nil
		}

		flagName, flagVal := matches[1], matches[2]
		remaining = matches[3]

		if flagName == "j" || flagName == "jump" {
			flagVal += remaining
			remaining = ""
		}
		opts[flagName] = flagVal
	}

	return util.MapAsString(opts)
}

func syncChain(ipt IPTables, table, chain string, target []string) error {
	curr, err := getRules(ipt, table, chain)
	if err != nil {
		return fmt.Errorf("iptables get: %s", err.Error())
	}

	_, rulesToDel, rulesToAdd := join.HashJoin(
		join.StringSlice(curr), join.StringSlice(target), ruleKey, ruleKey)

	for _, r := range rulesToDel {
		ruleSpec := strings.Split(r.(string), " ")
		if err := ipt.Delete(table, chain, ruleSpec...); err != nil {
			return fmt.Errorf("iptables delete: %s", err)
		}
	}

	for _, r := range rulesToAdd {
		ruleSpec := strings.Split(r.(string), " ")
		if err := ipt.Append(table, chain, ruleSpec...); err != nil {
			return fmt.Errorf("iptables append: %s", err)
		}
	}

	return nil
}

func getRules(ipt IPTables, table, chain string) (rules []string, err error) {
	rawRules, err := ipt.List(table, chain)
	if err != nil {
		return nil, err
	}

	for _, r := range rawRules {
		if !strings.HasPrefix(r, "-A") {
			continue
		}

		rSplit := strings.SplitN(r, " ", 3)
		if len(rSplit) != 3 {
			return nil, fmt.Errorf("malformed rule: %s", r)
		}

		rules = append(rules, rSplit[2])
	}

	return rules, nil
}

func preroutingRules(publicInterface string, containers []db.Container,
	connections []db.Connection) (rules []string) {

	// Map each label to all ports on which it can receive packets
	// from the public internet.
	portsFromWeb := make(map[string]map[int]struct{})
	for _, conn := range connections {
		if conn.From != stitch.PublicInternetLabel {
			continue
		}

		if _, ok := portsFromWeb[conn.To]; !ok {
			portsFromWeb[conn.To] = make(map[int]struct{})
		}

		portsFromWeb[conn.To][conn.MinPort] = struct{}{}
	}

	// Map the container's port to the same port of the host.
	for _, dbc := range containers {
		for _, label := range dbc.Labels {
			for port := range portsFromWeb[label] {
				for _, protocol := range []string{"tcp", "udp"} {
					rules = append(rules, fmt.Sprintf(
						"-i %[1]s -p %[2]s -m %[2]s "+
							"--dport %[3]d -j DNAT "+
							"--to-destination %[4]s:%[3]d",
						publicInterface, protocol, port, dbc.IP))
				}
			}
		}
	}

	return rules
}

func postroutingRules(publicInterface string, containers []db.Container,
	connections []db.Connection) (rules []string) {

	// Map each label to all ports on which it can send packets
	// to the public internet.
	portsToWeb := make(map[string]map[int]struct{})
	for _, conn := range connections {
		if conn.To != stitch.PublicInternetLabel {
			continue
		}

		if _, ok := portsToWeb[conn.From]; !ok {
			portsToWeb[conn.From] = make(map[int]struct{})
		}

		portsToWeb[conn.From][conn.MinPort] = struct{}{}
	}

	for _, dbc := range containers {
		for _, label := range dbc.Labels {
			for port := range portsToWeb[label] {
				for _, protocol := range []string{"tcp", "udp"} {
					rules = append(rules, fmt.Sprintf(
						"-s %[1]s/32 -p %[2]s -m %[2]s "+
							"--dport %[3]d -o %[4]s "+
							"-j MASQUERADE",
						dbc.IP, protocol, port, publicInterface,
					))
				}
			}
		}
	}

	return rules
}

type rule struct {
	table    string
	chain    string
	ruleSpec []string
}

func setDefaultRules(ipt IPTables) error {
	rules := []rule{
		{
			table:    "filter",
			chain:    "FORWARD",
			ruleSpec: []string{"-j", "ACCEPT"},
		},
		{
			table:    "nat",
			chain:    "INPUT",
			ruleSpec: []string{"-j", "ACCEPT"},
		},
		{
			table:    "nat",
			chain:    "OUTPUT",
			ruleSpec: []string{"-j", "ACCEPT"},
		},
	}
	for _, r := range rules {
		if err := ipt.AppendUnique(r.table, r.chain, r.ruleSpec...); err != nil {
			return fmt.Errorf("iptables append: %s", err)
		}
	}
	return nil
}

// getPublicInterfaceImpl gets the interface with the default route.
func getPublicInterfaceImpl() (string, error) {
	routes, err := routeList(nil, 0)
	if err != nil {
		return "", fmt.Errorf("route list: %s", err)
	}

	var defaultRoute *netlink.Route
	for _, r := range routes {
		if r.Dst == nil {
			defaultRoute = &r
			break
		}
	}

	if defaultRoute == nil {
		return "", errors.New("missing default route")
	}

	link, err := linkByIndex(defaultRoute.LinkIndex)
	if err != nil {
		return "", fmt.Errorf("default route missing interface: %s", err)
	}

	return link.Attrs().Name, err
}

var routeList = netlink.RouteList
var linkByIndex = netlink.LinkByIndex
var getPublicInterface = getPublicInterfaceImpl

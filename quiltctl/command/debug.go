package command

import (
	"archive/tar"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/quilt/quilt/api"
	"github.com/quilt/quilt/api/client"
	"github.com/quilt/quilt/api/client/getter"
	"github.com/quilt/quilt/db"
	"github.com/quilt/quilt/minion/supervisor/images"
	"github.com/quilt/quilt/quiltctl/ssh"
	"github.com/quilt/quilt/util"
)

// Stored in variables to be mocked out for the unit tests.
var timestamp = time.Now
var execCmd = exec.Command

const (
	containerDir = "containers"
	machineDir   = "machines"
)

// Debug contains the options for downloading debug logs from machines and containers.
type Debug struct {
	privateKey string
	outPath    string
	all        bool
	containers bool
	machines   bool
	tar        bool
	ids        []string

	common       *commonFlags
	clientGetter client.Getter
	sshGetter    ssh.Getter
}

type logTarget struct {
	ip       string
	dir      string
	stitchID string
	cmds     []logCmd
}

type logCmd struct {
	name string
	cmd  string
}

// A list of commands to run on the daemon. These must be formatted with the host address
// of the daemon. They will be prepended with 'quilt'.
var daemonCmds = []logCmd{
	{"ps", "ps -H=%s"},
	{"version", "version -H=%s"},
}

// A list of commands to output various machine logs.
var machineCmds = []logCmd{
	{"docker-ps", "docker ps -a"},
	{"minion", "docker logs minion"},
	{images.Etcd, "docker logs " + images.Etcd},
	{images.Ovncontroller, "docker logs " + images.Ovncontroller},
	{images.Ovnnorthd, "docker logs " + images.Ovnnorthd},
	{images.Ovsdb, "docker logs " + images.Ovsdb},
	{images.Ovsvswitchd, "docker logs " + images.Ovsvswitchd},
	{images.Registry, "docker logs " + images.Registry},
	{"syslog", "sudo cat /var/log/syslog"},
	{"journalctl", "sudo journalctl -xe"},
	{"uname", "uname"},
	{"dmesg", "dmesg"},
	{"uptime", "uptime"},
}

// A list of commands to output various container logs. Container commands need to be
// formatted with the DockerID.
var containerCmds = []logCmd{
	{"logs", "docker logs %s"},
}

// NewDebugCommand creates a new Debug command instance.
func NewDebugCommand() *Debug {
	return &Debug{
		clientGetter: getter.New(),
		sshGetter:    ssh.New,
		common:       &commonFlags{},
	}
}

var debugUsage = `usage: quilt debug-logs [-v] [-tar=<true/false>] [-i <keyfile>]` +
	` [-o=<path>] <-all | -containers | -machines | <id> ...>

Fetch logs for a set of machines or containers, placing
the contents in appropriately named file inside a
timestamped tarball. To store the contents in a folder
instead, use the flag '-tar=false'.

If the -all option is supplied, logs from all machines
and containers will be fetched. Else if the -containers
option is supplied, logs from all containers will be
fetched. The -machines option is similar. Otherwise if
none of the above options are given, a list of IDs can
be supplied, which may be a mix of machines and
containers in any order. Either one of these options
or a list of IDs must be supplied.

The -o option may be provided to optionally specify a
name for the tar file (or folder) instead of using a
timestamped name.

If -all is supplied, all other arguments are ignored. If
-containers or -machines are supplied, the list of IDs
is ignored, but they do not override each other. It follows
that the below commands are equivalent:
quilt debug-logs -all
quilt debug-logs -machines -containers
quilt debug-logs <supply all machine/container IDs>

To get the logs of machine 09ed35808a0b using a specific private key:
quilt debug-logs -i ~/.ssh/quilt 09ed35808a0b
`

// InstallFlags sets up parsing for command line flags.
func (dCmd *Debug) InstallFlags(flags *flag.FlagSet) {
	dCmd.common.InstallFlags(flags)
	flags.StringVar(&dCmd.privateKey, "i", "",
		"the private key to use to connect to the host")
	flags.StringVar(&dCmd.outPath, "o", "",
		"output path for the logs (defaults to timestamped path)")
	flags.BoolVar(&dCmd.all, "all", false, "if provided, fetch all debug logs")
	flags.BoolVar(&dCmd.containers, "containers", false,
		"if provided, fetch all debug logs for application containers")
	flags.BoolVar(&dCmd.machines, "machines", false,
		"if provided, fetch all debug logs for machines"+
			" (including quilt system containers)")
	flags.BoolVar(&dCmd.tar, "tar", true,
		"if true (default), compress the logs into a tarball")

	flags.Usage = func() {
		fmt.Println(debugUsage)
		flags.PrintDefaults()
	}
}

// Parse parses the command line arguments for the debug command.
func (dCmd *Debug) Parse(args []string) error {
	dCmd.ids = args
	if len(dCmd.ids) == 0 && !(dCmd.all || dCmd.machines || dCmd.containers) {
		return errors.New("must supply at least one ID or set option")
	}

	return nil
}

// Run downloads debug logs from the relevant machines and containers.
func (dCmd Debug) Run() int {
	if dCmd.outPath == "" {
		dCmd.outPath = fmt.Sprintf("debug_logs_%s",
			timestamp().Format("Mon_Jan_02_15-04-05"))
	}
	if err := util.Mkdir(dCmd.outPath, 0755); err != nil {
		log.Error(err)
		return 1
	}

	c, err := dCmd.clientGetter.Client(dCmd.common.host)
	if err != nil {
		log.Error(err)
		return 1
	}
	defer c.Close()

	machines, err := c.QueryMachines()
	if err != nil {
		log.Error(err)
		return 1
	}

	ipMap := map[string]string{}
	containers := []db.Container{}
	for _, m := range machines {
		ipMap[m.PrivateIP] = m.PublicIP
		if m.Role != db.Worker || m.PublicIP == "" {
			continue
		}

		cli, err := dCmd.clientGetter.Client(api.RemoteAddress(m.PublicIP))
		if err != nil {
			log.Error(err)
			return 1
		}
		defer cli.Close()

		conts, err := cli.QueryContainers()
		if err != nil {
			log.Error(err)
			return 1
		}

		containers = append(containers, conts...)
	}

	dCmd.machines = dCmd.machines || dCmd.all
	dCmd.containers = dCmd.containers || dCmd.all

	var targets []logTarget
	mTargets := machinesToTargets(machines)
	cTargets := containersToTargets(containers, ipMap)
	if !(dCmd.machines || dCmd.containers) {
		targets = append(append(targets, cTargets...), mTargets...)
		if targets, err = filterTargets(targets, dCmd.ids); err != nil {
			log.Error(err)
			return 1
		}
	}

	if dCmd.machines {
		targets = append(targets, mTargets...)
	}
	if dCmd.containers {
		targets = append(targets, cTargets...)
	}

	return dCmd.downloadLogs(targets)
}

func (dCmd Debug) downloadLogs(targets []logTarget) int {
	rootDir := dCmd.outPath
	if err := util.Mkdir(filepath.Join(rootDir, machineDir), 0755); err != nil {
		log.Error(err)
		return 1
	}

	if err := util.Mkdir(filepath.Join(rootDir, containerDir), 0755); err != nil {
		log.Error(err)
		return 1
	}

	// Since we don't want the failure of downloading one or more logs to affect the
	// rest, errors that arise from the fetching of logs are ignored and errno is
	// simply incremented. The debug-logs command's exit code is errno if this line
	// of the code is reached before exiting.
	var errno int
	for _, cmd := range daemonCmds {
		file := filepath.Join(rootDir, cmd.name)
		fmtCmd := fmt.Sprintf(cmd.cmd, dCmd.common.host)
		qCmd := execCmd("quilt", strings.Fields(fmtCmd)...)
		log.Debugf("Gathering `quilt %s` output", fmtCmd)
		if result, err := qCmd.CombinedOutput(); err != nil {
			errno++
			log.Error(err)
		} else if err := util.WriteFile(file, result, 0644); err != nil {
			errno++
			log.Error(err)
		}
	}

	for _, t := range targets {
		path := filepath.Join(rootDir, t.dir, t.stitchID)
		if err := util.Mkdir(path, 0755); err != nil {
			errno++
			log.Error(err)
			continue
		}

		conn, err := dCmd.sshGetter(t.ip, dCmd.privateKey)
		if err != nil {
			errno++
			log.Error(err)
			continue
		}

		for _, cmd := range t.cmds {
			log.Debugf("Downloading log '%s' for target %s", cmd.name,
				t.stitchID)

			result, err := conn.CombinedOutput(cmd.cmd)
			if err != nil {
				// Since workers and masters have different containers,
				// some downloads will fail and this is normal. Because
				// of this (and to keep the code simpler), errors while
				// fetching logs are at the debug level and errno is not
				// incremented. Connection problems will be caught by the
				// attempt to establish an ssh connection to the host
				// machine.
				log.WithError(err).Debugf("Failed to get log '%s'"+
					" from target %s", cmd.name, t.stitchID)
				continue
			}

			logFile := filepath.Join(path, cmd.name)
			if err := util.WriteFile(logFile, result, 0644); err != nil {
				errno++
				log.Error(err)
			}
		}
	}

	if errno > 0 {
		log.Error("Some downloads failed")
	}

	if dCmd.tar {
		errno += tarball(rootDir)
	}
	return errno
}

func machinesToTargets(machines []db.Machine) []logTarget {
	targets := []logTarget{}
	for _, m := range machines {
		if m.PublicIP == "" {
			continue
		}

		t := logTarget{
			ip:       m.PublicIP,
			dir:      machineDir,
			stitchID: m.StitchID,
			cmds:     nil,
		}
		t.cmds = append(t.cmds, machineCmds...)
		targets = append(targets, t)
	}
	return targets
}

func containersToTargets(containers []db.Container, ips map[string]string) []logTarget {
	targets := []logTarget{}
	for _, c := range containers {
		if c.Minion == "" {
			continue
		}

		ip, ok := ips[c.Minion]
		if !ok {
			log.Errorf("No machine with private IP %s", c.Minion)
			continue
		}

		t := logTarget{
			ip:       ip,
			dir:      containerDir,
			stitchID: c.StitchID,
			cmds:     nil,
		}
		for _, cmd := range containerCmds {
			cmd.cmd = fmt.Sprintf(cmd.cmd, c.DockerID)
			t.cmds = append(t.cmds, cmd)
		}
		targets = append(targets, t)
	}
	return targets
}

func filterTargets(targets []logTarget, ids []string) ([]logTarget, error) {
	result := []logTarget{}
	for _, id := range ids {
		t, err := findTarget(targets, id)
		if err != nil {
			return result, err
		}

		result = append(result, t)
	}
	return result, nil
}

func findTarget(targets []logTarget, id string) (logTarget, error) {
	choice := logTarget{stitchID: ""}
	for _, t := range targets {
		if len(id) > len(t.stitchID) || t.stitchID[:len(id)] != id {
			continue
		}

		if choice.stitchID != "" {
			return logTarget{}, fmt.Errorf("ambiguous stitchIDs %s and %s",
				choice.stitchID, t.stitchID)
		}
		choice = t
	}

	if choice.stitchID == "" {
		return logTarget{}, fmt.Errorf("no target with stitchID %s", id)
	}

	return choice, nil
}

func tarball(root string) int {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	ballUp := func(path string, info os.FileInfo, err error) error {
		hdr, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		hdr.Name = path

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		contents, err := util.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = tw.Write([]byte(contents))
		return err
	}

	if err := util.Walk(root, ballUp); err != nil {
		log.WithError(err).Errorf("Failed to tarball directory %s", root)
		return 1
	}

	if err := tw.Close(); err != nil {
		log.WithError(err).Error("Failed to close tar writer")
		return 1
	}

	if err := util.RemoveAll(root); err != nil {
		log.WithError(err).Error("Failed to remove temporary log directory")
		return 1
	}

	if err := util.WriteFile(root+".tar", buf.Bytes(), 0644); err != nil {
		log.WithError(err).Error("Failed to write tarball")
		return 1
	}

	return 0
}

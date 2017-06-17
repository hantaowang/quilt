package inspect

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/quilt/quilt/stitch"
	"github.com/quilt/quilt/util"
)

func getSlug(configPath string) (string, error) {
	var slug string
	for i, ch := range configPath {
		if ch == '.' {
			slug = configPath[:i]
			break
		}
	}
	if len(slug) == 0 {
		return "", fmt.Errorf("could not find proper output file name")
	}

	return slug, nil
}

func viz(configPath string, blueprint stitch.Stitch, graph stitch.Graph,
	outputFormat string) {
	slug, err := getSlug(configPath)
	if err != nil {
		panic(err)
	}
	dot := makeGraphviz(graph)
	graphviz(outputFormat, slug, dot)
}

func makeGraphviz(graph stitch.Graph) string {
	dotfile := "strict digraph {\n"

	for i, av := range graph.Availability {
		dotfile += subGraph(i, av.Nodes()...)
	}

	var lines []string
	for _, edge := range graph.GetConnections() {
		lines = append(lines,
			fmt.Sprintf(
				"    %s -> %s\n",
				edge.From,
				edge.To,
			),
		)
	}

	sort.Strings(lines)
	for _, line := range lines {
		dotfile += line + "\n"
	}

	dotfile += "}\n"

	return dotfile
}

func subGraph(i int, labels ...string) string {
	subgraph := fmt.Sprintf("    subgraph cluster_%d {\n", i)
	str := ""
	sort.Strings(labels)
	for _, l := range labels {
		str += l + "; "
	}
	subgraph += "        " + str + "\n    }\n"
	return subgraph
}

// Graphviz generates a specification for the graphviz program that visualizes the
// communication graph of a stitch.
func graphviz(outputFormat string, slug string, dot string) {
	f, err := util.AppFs.Create(slug + ".dot")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write([]byte(dot))
	if outputFormat == "graphviz" {
		return
	}
	defer exec.Command("rm", slug+".dot").Run()

	// Dependencies:
	// - easy-graph (install Graph::Easy from cpan)
	// - graphviz (install from your favorite package manager)
	var writeGraph *exec.Cmd
	switch outputFormat {
	case "ascii":
		writeGraph = exec.Command("graph-easy", "--input="+slug+".dot",
			"--as_ascii")
	case "pdf":
		writeGraph = exec.Command("dot", "-Tpdf", "-o", slug+".pdf",
			slug+".dot")
	}
	writeGraph.Stdout = os.Stdout
	writeGraph.Stderr = os.Stderr
	writeGraph.Run()
}

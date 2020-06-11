package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kechako/go-yahoo/da"
	"github.com/kechako/mktree/node"
)

var appID string

func init() {
	flag.StringVar(&appID, "appid", os.Getenv("YAHOO_APP_ID"), "Yahoo! Application ID.")
}

func printHelp() {
	flag.Usage()
	flag.PrintDefaults()
}

func run() (int, error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s text\n\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 1 {
		printHelp()
		return 2, nil
	}

	if appID == "" {
		fmt.Fprint(os.Stderr, "Yahoo! Application ID is not specified.\n\n")
		printHelp()
		return 2, nil
	}

	client := da.New(appID)
	res, err := client.Parse(context.Background(), flag.Arg(0))
	if err != nil {
		return 1, fmt.Errorf("could not parse the text: %w", err)
	}

	depMap := make(map[int][]da.Chunk)
	for _, chunk := range res.Chunks {
		depMap[chunk.Head] = append(depMap[chunk.Head], chunk)
	}

	root := makeNode(depMap[-1][0], depMap, 0)

	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()

	root.Print(buf)

	return 0, nil
}

func makeNode(chunk da.Chunk, depMap map[int][]da.Chunk, indent int) *node.Node {
	n := node.NewIndent(chunk.String(), indent)

	for _, dep := range depMap[chunk.ID] {
		child := makeNode(dep, depMap, indent+1)
		n.Add(child)
	}

	return n
}

func main() {
	code, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v\n", err)
	}
	os.Exit(code)
}

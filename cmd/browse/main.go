package main

import (
	"fmt"
	"os"

	"github.com/deglebe/browse/pkg/html"
)

// print dom tree of passed file
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: browse <path-to-html>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()

	parser := html.NewParser(f)
	domTree, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Parse error:", err)
		os.Exit(1)
	}

	domTree.PrettyPrint(0)
}

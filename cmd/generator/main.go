package main

import (
	"fmt"
	"os"

	"github.com/grundleborg/mattermost-model-generator/internal/model"
	"github.com/grundleborg/mattermost-model-generator/internal/parser"
	"github.com/grundleborg/mattermost-model-generator/internal/renderer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "not enough program arguments")
		os.Exit(1)
	}
	// TODO: take output directory as an arg/flag
	// TODO: take packages to generate as an arg/flag
	_ = os.Mkdir("output", 0755)
	for _, arg := range os.Args[1:] {
		str, err := parser.ParseFile(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse: %s\n", err)
			os.Exit(1)
		}

		if err := renderer.RenderToFile("output", model.Model, str); err != nil {
			fmt.Fprintf(os.Stderr, "could not render %q to file: %s\n", str.Type, err)
			os.Exit(2)
		}
		if err := renderer.RenderToFile("output", model.Client, str); err != nil {
			fmt.Fprintf(os.Stderr, "could not render %q to file: %s\n", str.Type, err)
			os.Exit(2)
		}
	}
}

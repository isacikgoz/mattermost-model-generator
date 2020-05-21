package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/grundleborg/mattermost-model-generator/internal/finder"
)

// go run ./cmd/migrator $GOPATH/src/github.com/mattermost/mattermost-server/**/*.go
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "not enough program arguments")
		os.Exit(1)
	}

	fs := token.NewFileSet()

	for _, arg := range os.Args[1:] {
		f, err := parser.ParseFile(fs, arg, nil, parser.AllErrors)
		if err != nil {
			log.Printf("could not parse %s: %v", arg, err)
			continue
		}
		w := &finder.Walker{
			Name:    "Post",
			Package: "model",
		}
		ast.Walk(w, f)
	}
}

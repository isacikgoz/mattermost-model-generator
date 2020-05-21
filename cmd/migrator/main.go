package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/grundleborg/mattermost-model-generator/internal/finder"
)

// go run ./cmd/migrator model Post $GOPATH/src/github.com/mattermost/mattermost-server/**/*.go
func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "not enough program arguments: [cmd] [package] [model] [files to scan]")
		os.Exit(1)
	}
	packageName := os.Args[1]
	modelName := os.Args[2]
	fset := token.NewFileSet()

	for _, fileName := range os.Args[3:] {
		if strings.Contains(fileName, "mattermost-server/vendor") {
			continue
		}
		fileNode, err := parser.ParseFile(fset, fileName, nil, parser.AllErrors|parser.ParseComments)
		if err != nil {
			log.Printf("could not parse %s: %v", fileName, err)
			continue
		}
		w := &finder.Walker{
			Name:    modelName,
			Package: packageName,
		}
		ast.Walk(w, fileNode)
		buf := new(bytes.Buffer)
		if err = format.Node(buf, fset, fileNode); err != nil {
			log.Printf("error: %v\n", err)
		} else if err := ioutil.WriteFile(fileName, buf.Bytes(), 0664); err != nil {
			log.Printf("error: %v\n", err)
		}
	}
}

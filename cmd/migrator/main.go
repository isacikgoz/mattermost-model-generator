package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/grundleborg/mattermost-model-generator/internal/finder"
)

// go run ./cmd/migrator model Post $GOPATH/src/github.com/mattermost/mattermost-server
func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "not enough program arguments: [cmd] [package] [model] [mattermost-server dir]")
		os.Exit(1)
	}
	packageName := os.Args[1]
	modelName := os.Args[2]
	var goFiles []string
	walker := func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if (info.IsDir() && strings.HasPrefix(info.Name(), ".")) || strings.Contains(fullPath, "mattermost-server/vendor") {
			return nil
		}
		if path.Ext(info.Name()) == ".go" && !strings.HasSuffix(info.Name(), "_test") {
			goFiles = append(goFiles, fullPath)
		}
		return nil
	}
	if err := filepath.Walk(os.Args[3], walker); err != nil {
		log.Fatalf("could not scan folder: %v", err)
	}
	fset := token.NewFileSet()

	for _, fileName := range goFiles {
		fileNode, err := parser.ParseFile(fset, fileName, nil, parser.AllErrors|parser.ParseComments)
		if err != nil {
			log.Printf("could not parse %s: %v", fileName, err)
			continue
		}
		w := &finder.Walker{
			Name:    modelName,
			Package: packageName,
		}
		result := w.Process(fileNode)
		buf := new(bytes.Buffer)
		if err = format.Node(buf, fset, result); err != nil {
			log.Printf("error: %v\n", err)
		} else if err := ioutil.WriteFile(fileName, buf.Bytes(), 0664); err != nil {
			log.Printf("error: %v\n", err)
		}
	}
}

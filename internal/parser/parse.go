package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/fatih/structtag"
	"github.com/grundleborg/mattermost-model-generator/internal/model"
)

// ParseFile reads a file and generates representation of structs to be generated.
func ParseFile(path string) ([]*model.Struct, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("could not parse file %q: %w", path, err)
	}

	var structs []*model.Struct
	ast.Inspect(file, func(n ast.Node) bool {
		tp, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		st, ok := tp.Type.(*ast.StructType)
		if !ok {
			return true
		}
		fields, err := parseFields(st.Fields.List)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse fields for %q: %s", tp.Name.Name, err)
			return false // maybe replace this with a panic?
		}
		structs = append(structs, &model.Struct{
			Type:   tp.Name.Name,
			Fields: fields,
		})
		return true
	})

	return structs, nil
}

func parseFields(fields []*ast.Field) ([]*model.Field, error) {
	var fs []*model.Field
	for _, field := range fields {
		st, err := structtag.Parse(strings.ReplaceAll(field.Tag.Value, "`", ""))
		if err != nil {
			return nil, err
		}

		tags := make(map[string][]string, len(st.Tags()))
		for _, tag := range st.Tags() {
			tags[tag.Key] = append([]string{tag.Name}, tag.Options...)
		}
		fs = append(fs, &model.Field{
			Name: field.Names[0].Name,
			Type: field.Type.(*ast.Ident).Name,
			Tags: tags,
		})
	}
	return fs, nil
}

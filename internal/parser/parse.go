package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/fatih/structtag"
	"github.com/grundleborg/mattermost-model-generator/internal/model"
)

type File struct {
	Struct *model.Struct       // the main type for the file
	Types  []*model.CustomType // custom types for this struct
}

func FormatNode(node ast.Node) string {
	buf := new(bytes.Buffer)
	_ = format.Node(buf, token.NewFileSet(), node)
	return buf.String()
}

// ParseFile reads a file and generates representation of structs to be generated.
func ParseFile(path string) (*model.Struct, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("could not parse file %q: %w", path, err)
	}

	var str *model.Struct
	var types []*model.CustomType
	ast.Inspect(file, func(n ast.Node) bool {
		tp, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		switch t := tp.Type.(type) {
		case *ast.StructType:
			fields, err := parseFields(t.Fields.List)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not parse fields for %q: %s", tp.Name.Name, err)
				return false // maybe replace this with a panic?
			}
			str = &model.Struct{
				Type:   tp.Name.Name,
				Fields: fields,
			}
			return true
		case *ast.Ident:
			types = append(types, &model.CustomType{
				Name:           tp.Name.Name,
				UnderlyingType: t.Name,
			})
			return true
		}
		return true
	})

	str.CustomTypes = types
	return str, nil
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
			Type: FormatNode(field.Type),
			Tags: tags,
		})
	}
	return fs, nil
}

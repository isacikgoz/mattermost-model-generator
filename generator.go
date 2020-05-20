// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/fatih/structtag"
)

type ModelParams struct {
	Type   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
	Tags map[string][]string
}

func main() {
	fmt.Println("Starting code generation.")
	processFile("model.go")
	fmt.Println("Code generation completed.")
}

func initTemplate(name, file string) *template.Template {
	data, err := ioutil.ReadFile("templates/" + file)
	if err != nil {
		panic(err)
	}
	funcMap := template.FuncMap{
		// make s string start with upper case
		"public": func(s string) string {
			return strings.Title(s)
		},
		// make s string start with upper case
		"receiver": func(s string) string {
			return strings.ToLower(string(s[0]))
		},
		// prints only json tags for a field
		"json": func(tags map[string][]string) string {
			for k, v := range tags {
				if k == "json" {
					return "`json:\"" + strings.Join(v, ",") + "\"`"
				}
			}
			return ""
		},
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(string(data))
	if err != nil {
		panic(err)
	}

	return tmpl
}

func processFile(filePath string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		tp, ok := n.(*ast.TypeSpec)
		if ok {
			typeName := tp.Name.Name
			st, ok := tp.Type.(*ast.StructType)
			if ok {
				processStruct(st, typeName)
			}

		}
		return true
	})
}

func processStruct(st *ast.StructType, typeName string) {
	params := ModelParams{
		Type:   typeName,
		Fields: generateFields(st.Fields.List),
	}
	for _, pkg := range []string{"client", "model"} {
		buf := new(bytes.Buffer)
		buf.Write(renderFile(pkg, params))
		ioutil.WriteFile(pkg+"/"+strings.ToLower(typeName)+".go", buf.Bytes(), 0664)
	}
}

func generateFields(fields []*ast.Field) []Field {
	var fs []Field
	for _, field := range fields {
		st, err := structtag.Parse(strings.ReplaceAll(field.Tag.Value, "`", ""))
		if err != nil {
			panic(err)
		}

		tags := make(map[string][]string)
		for _, tag := range st.Tags() {
			tags[tag.Key] = append([]string{tag.Name}, tag.Options...)
		}
		fs = append(fs, Field{
			Name: field.Names[0].Name,
			Type: field.Type.(*ast.Ident).Name,
			Tags: tags,
		})
	}
	return fs
}

func renderFile(pkg string, params ModelParams) []byte {
	buf := new(bytes.Buffer)
	tmpl := initTemplate(pkg, pkg+".go.tmpl")
	err := tmpl.Execute(buf, params)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

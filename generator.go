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
)

type Templates struct {
	File   *template.Template
	Getter *template.Template
	Initializer *template.Template
}

var templates Templates

type FileParams struct {
	Getters string
	Initializer string
}

type GetterParams struct {
	ReceiverType string
	ReceiverName string
	FieldType    string
	FieldName    string
	FuncName     string
}

type InitializerParams struct {
	StructName string
	ParamName string
	Fields []InitializerField
}

type InitializerField struct {
	Name string
	InternalName string
	JSONName string
	Type string
}

func main() {
	fmt.Println("Starting code generation.")
	initTemplates()
	processFile("model.go")
	fmt.Println("Code generation completed.")
}

func initTemplates() {
	templates.File = initTemplate("file", "file.go.tmpl")
	templates.Getter = initTemplate("getter", "getter.go.tmpl")
	templates.Initializer = initTemplate("initializer", "initializer.go.tmpl")
}

func initTemplate(name, file string) *template.Template {
	data, err := ioutil.ReadFile("templates/"+file)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New(name).Parse(string(data))
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
	getters := processGetters(st.Fields.List, typeName)
	initializer := processInitializer(st.Fields.List, typeName)

	params := FileParams{
		Getters: getters,
		Initializer: initializer,
	}

	buf := new(bytes.Buffer)
	buf.Write(generateFile(params))
	ioutil.WriteFile("model/"+strings.ToLower(typeName)+".go", buf.Bytes(), 0664)
}

func processGetters(fields []*ast.Field, typeName string) string {
	gettersBuf := new(bytes.Buffer)

	for _, field := range fields {
		fieldName := field.Names[0].Name
		fieldType := field.Type.(*ast.Ident).Name
		gettersBuf.Write(generateGetter(typeName, fieldName, fieldType))
	}

	return string(gettersBuf.Bytes())
}

func processInitializer(fields []*ast.Field, typeName string) string {
	buf := new(bytes.Buffer)
	initializerParams := InitializerParams{
		StructName: typeName,
		ParamName: strings.ToLower(string(typeName[0])),
	}

	for _, field := range fields {
		fieldName := field.Names[0].Name
		fieldType := field.Type.(*ast.Ident).Name
		// FIXME: Parse the struct tags properly.
		jsonName := field.Tag.Value
		initializerField := InitializerField{
			Name: strings.ToUpper(string(fieldName[0])) + fieldName[1:],
			InternalName: fieldName,
			JSONName: jsonName,
			Type: fieldType,
		}
		initializerParams.Fields = append(initializerParams.Fields, initializerField)
	}

	buf.Write(generateInitializer(initializerParams))
	return string(buf.Bytes())
}

func generateFile(params FileParams) []byte {
	buf := new(bytes.Buffer)
	err := templates.File.Execute(buf, params)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func generateGetter(receiverType, fieldName, fieldType string) []byte {
	params := GetterParams{
		ReceiverType: receiverType,
		ReceiverName: strings.ToLower(string(receiverType[0])),
		FieldType:    fieldType,
		FieldName:    fieldName,
		FuncName:     strings.ToUpper(string(fieldName[0])) + fieldName[1:],
	}

	buf := new(bytes.Buffer)
	err := templates.Getter.Execute(buf, params)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func generateInitializer(params InitializerParams) []byte {
	buf := new(bytes.Buffer)
	err := templates.Initializer.Execute(buf, params)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

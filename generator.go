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

type ModelParams struct {
	BaseTypeName string
	ReceiverName string
	Getters      []GetterParams
	Initializer  InitializerParams
}

type GetterParams struct {
	ReceiverType string
	ReceiverName string
	FieldType    string
	FieldName    string
	FuncName     string
}

type InitializerParams struct {
	TypeName string
	Fields   []InitializerField
}

type InitializerField struct {
	Name         string
	InternalName string
	JSONTag      string
	Type         string
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
	params := ModelParams{
		BaseTypeName: typeName,
		ReceiverName: makeReceiverName(typeName),
		Getters:      generateGetters(st.Fields.List, typeName),
		Initializer:  generateInitializer(st.Fields.List, typeName),
	}

	buf := new(bytes.Buffer)
	buf.Write(renderFile(params))
	ioutil.WriteFile("model/"+strings.ToLower(typeName)+".go", buf.Bytes(), 0664)
}

func generateGetters(fields []*ast.Field, typeName string) []GetterParams {
	params := []GetterParams{}
	for _, field := range fields {
		fieldName := field.Names[0].Name
		fieldType := field.Type.(*ast.Ident).Name

		params = append(params, GetterParams{
			ReceiverType: typeName,
			ReceiverName: makeReceiverName(typeName),
			FieldType:    fieldType,
			FieldName:    fieldName,
			FuncName:     strings.ToUpper(string(fieldName[0])) + fieldName[1:],
		})
	}

	return params
}

func generateInitializer(fields []*ast.Field, baseTypeName string) InitializerParams {
	params := InitializerParams{
		TypeName: fmt.Sprintf("%sInitializer", baseTypeName),
	}

	for _, field := range fields {
		fieldName := field.Names[0].Name
		fieldType := field.Type.(*ast.Ident).Name
		// FIXME: Parse the struct tags properly.
		jsonTag := field.Tag.Value

		params.Fields = append(params.Fields, InitializerField{
			Name:         strings.ToUpper(string(fieldName[0])) + fieldName[1:],
			InternalName: fieldName,
			JSONTag:      jsonTag,
			Type:         fieldType,
		})
	}

	return params
}

func renderFile(params ModelParams) []byte {
	buf := new(bytes.Buffer)
	tmpl := initTemplate("file", "file.go.tmpl")
	err := tmpl.Execute(buf, params)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func makeReceiverName(receiverType string) string {
	return strings.ToLower(string(receiverType[0]))
}

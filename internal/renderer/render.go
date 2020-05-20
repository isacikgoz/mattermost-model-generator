package renderer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/grundleborg/mattermost-model-generator/internal/model"
)

// Render generates the file for the struct
func Render(pkg model.Package, st *model.Struct) ([]byte, error) {
	tmpl, err := initTemplate(string(pkg))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, st)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func RenderToFile(dir string, pkg model.Package, st *model.Struct) error {
	buf := new(bytes.Buffer)
	bytes, err := Render(pkg, st)
	if err != nil {
		return err
	}
	buf.Write(bytes)

	os.MkdirAll(filepath.Join(dir, string(pkg)), 0755)
	fmt.Println(filepath.Join(dir, string(pkg)))

	dstFile := filepath.Join(dir, string(pkg), strings.ToLower(st.Type)+".go")
	return ioutil.WriteFile(dstFile, buf.Bytes(), 0664)
}

func initTemplate(name string) (*template.Template, error) {
	data, err := ioutil.ReadFile("templates/" + name + ".go.tmpl")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(string(data))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

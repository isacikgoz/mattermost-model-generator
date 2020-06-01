package renderer

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/grundleborg/mattermost-model-generator/internal/model"
)

var funcMap = template.FuncMap{
	// make s string start with upper case
	"public": func(s string) string {
		return strings.Title(s)
	},
	"sliceTypes": func(s []*model.Field) map[string]string {
		ret := map[string]string{}
		for _, r := range s {
			if strings.Contains(r.Type, "[]") {
				ret["SliceOf"+strings.Title(r.Type[2:])] = r.Type[2:]
			}
		}
		return ret
	},
	"generateInitializer": func(t, name, ft string) string {
		s := fmt.Sprintf("%s.%s", strings.ToLower(string(t[0])), strings.Title(name))
		if strings.Contains(ft, "[]") {
			return fmt.Sprintf("NewSliceOf%s(%s)", strings.Title(ft[2:]), s)
		}
		return s
	},
	"generateSetStatement": func(t string) string {
		if strings.Contains(t, "[]") {
			return ".Replace(v)"
		}
		return " = v"
	},
	"generateGetStatement": func(t string) string {
		if strings.Contains(t, "[]") {
			return ".Range()"
		}
		return ""
	},
	"processType": func(s string) string {
		if strings.Contains(s, "[]") {
			return "SliceOf" + strings.Title(s[2:])
		}
		return s
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
	// return true if any patch tag exists for struct fields
	"include_patch": func(fields []*model.Field) bool {
		for _, f := range fields {
			if patch_field(f.Tags) {
				return true
			}
		}
		return false
	},
	// return true if model:patch tag exists
	"patch_field": patch_field,
}

func patch_field(tags map[string][]string) bool {
	for k, v := range tags {
		if k == "model" {
			for _, s := range v {
				if s == "patch" {
					return true
				}
			}
		}
	}
	return false
}

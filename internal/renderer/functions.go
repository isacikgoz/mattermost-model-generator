package renderer

import (
	"strings"
	"text/template"

	"github.com/grundleborg/mattermost-model-generator/internal/model"
)

var funcMap = template.FuncMap{
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

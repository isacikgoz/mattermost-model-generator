package renderer

import (
	"strings"
	"text/template"
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
}

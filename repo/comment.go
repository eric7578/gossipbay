package repo

import (
	"text/template"
)

var (
	mdTmpl *template.Template
)

func init() {
	var (
		err       error
		mdComment = `{{ range . }}
{{ range .Posts }}[{{ .Title }}]({{ .URL }}) **{{ .NumUp }}** 推 **{{ .NumDown }}** 噓
{{end}}
{{end}}
`
	)

	mdTmpl, err = template.New("mdComment").Parse(mdComment)
	if err != nil {
		panic(err)
	}
}

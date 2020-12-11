package ntemplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteTemplateHTML(t *testing.T) {
	a := assert.New(t)
	t1, err := NewHTMLTemplate("t1", `{{- define "T1"}}ONE{{end}}
	{{- define "T2"}}<a href="/search?q={{.}}">{{.}}</a>{{end}}
	{{- define "T3"}}{{template "T1"}} {{template "T2"}}{{end}}
	{{- template "T3"}}`)
	a.Nil(err)
	str, err := t1.Execute(nil)
	a.Nil(err)
	a.Equal(`ONE <a href="/search?q="></a>`, str)
	str, err = t1.ExecuteTemplate("T2", "<b>a</b>")
	a.Nil(err)
	a.Equal(`<a href="/search?q=%3cb%3ea%3c%2fb%3e">&lt;b&gt;a&lt;/b&gt;</a>`, str)
}

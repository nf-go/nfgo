// Copyright 2020 The nfgo Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ntemplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustParseHtmlTemplate(t *testing.T) {
	t.Log(tmpl)
	htmlTmpl := MustParseHTMLTemplate(tmpl, "*.go")
	str, err := htmlTmpl.ExecuteTemplate("html.go", nil)
	a := assert.New(t)
	a.Nil(err)
	a.Contains(str, "type HTMLTemplate struct {")
	tt := htmlTmpl.Lookup("html.go")
	a.NotNil(tt)
	str, err = tt.Execute(nil)
	a.Nil(err)
	a.Contains(str, "type HTMLTemplate struct {")
	tt = htmlTmpl.Lookup("notexist.go")
	a.Nil(tt)

	str, err = htmlTmpl.ExecuteTemplate("text.go", nil)
	a.Nil(err)
	a.Contains(str, "type TextTemplate struct {")
}

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

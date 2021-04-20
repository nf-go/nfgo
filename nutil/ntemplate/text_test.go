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
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed html.go text.go
var tmpl embed.FS

func TestMustParseTextTemplate(t *testing.T) {
	t.Log(tmpl)
	textTmpl := MustParseTextTemplate(tmpl, "*.go")
	str, err := textTmpl.ExecuteTemplate("html.go", nil)
	a := assert.New(t)
	a.Nil(err)
	a.Contains(str, "type HTMLTemplate struct {")
	tt := textTmpl.Lookup("html.go")
	a.NotNil(tt)
	str, err = tt.Execute(nil)
	a.Nil(err)
	a.Contains(str, "type HTMLTemplate struct {")
	tt = textTmpl.Lookup("notexist.go")
	a.Nil(tt)

	str, err = textTmpl.ExecuteTemplate("text.go", nil)
	a.Nil(err)
	a.Contains(str, "type TextTemplate struct {")
}

func TestNewTextTemplate(t *testing.T) {
	a := assert.New(t)
	t1, err := NewTextTemplate("t1", `hello{{  . }}`)
	a.Nil(err)
	a.NotNil(t1)
	t2, err := NewTextTemplate("t2", `hello{{ foo  . }}`)
	a.NotNil(err)
	a.Nil(t2)
	t3 := MustNewTextTemplate("t3", `{{ . }} world`)
	a.NotNil(t3)
}

func TestExecute(t *testing.T) {
	a := assert.New(t)
	t1, err := NewTextTemplate("t1", `hello {{  . }}`)
	a.Nil(err)
	str, err := t1.Execute("world")
	a.Nil(err)
	a.Equal("hello world", str)
}

func TestExecuteTemplate(t *testing.T) {
	a := assert.New(t)
	t1, err := NewTextTemplate("t1", `{{- define "T1"}}ONE{{end}}
	{{- define "T2"}}TWO{{end}}
	{{- define "T3"}}{{template "T1"}} {{template "T2"}}{{end}}
	{{- template "T3"}}`)
	a.Nil(err)
	str, err := t1.Execute(nil)
	a.Nil(err)
	a.Equal("ONE TWO", str)
	str, err = t1.ExecuteTemplate("T2", nil)
	a.Nil(err)
	a.Equal("TWO", str)
}

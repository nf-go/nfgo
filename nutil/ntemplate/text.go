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
	"io/fs"
	"strings"
	"text/template"

	"nfgo.ga/nfgo/nlog"
)

// TextTemplate -
type TextTemplate struct {
	tmpl *template.Template
}

// ParseTextTemplate -
func ParseTextTemplate(fs fs.FS, patterns ...string) (*TextTemplate, error) {
	tmpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return nil, err
	}
	return &TextTemplate{
		tmpl: tmpl,
	}, nil
}

// MustParseTextTemplate -
func MustParseTextTemplate(fs fs.FS, patterns ...string) *TextTemplate {
	t, err := ParseTextTemplate(fs, patterns...)
	if err != nil {
		nlog.Fatal("fail to parse text template: ", err)
	}
	return t
}

// NewTextTemplate -
func NewTextTemplate(name, text string) (*TextTemplate, error) {
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return nil, err
	}
	return &TextTemplate{
		tmpl: tmpl,
	}, nil
}

// MustNewTextTemplate -
func MustNewTextTemplate(name, text string) *TextTemplate {
	t, err := NewTextTemplate(name, text)
	if err != nil {
		nlog.Fatal("fail to create text template: ", err)
	}
	return t
}

// Lookup - returns the template with the given name that is associated with t.
// It returns nil if there is no such template or the template has no definition.
func (t *TextTemplate) Lookup(name string) *TextTemplate {
	tmpl := t.tmpl.Lookup(name)
	if tmpl == nil {
		return nil
	}
	return &TextTemplate{
		tmpl: tmpl,
	}
}

// Execute - applies the template to the specified data object
func (t *TextTemplate) Execute(data interface{}) (string, error) {
	var sb strings.Builder
	if err := t.tmpl.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// ExecuteTemplate - applies the template associated with t that has the given name to the specified data object
func (t *TextTemplate) ExecuteTemplate(name string, data interface{}) (string, error) {
	var sb strings.Builder
	if err := t.tmpl.ExecuteTemplate(&sb, name, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

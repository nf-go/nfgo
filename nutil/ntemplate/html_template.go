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
	"html/template"
	"strings"

	"nfgo.ga/nfgo/nlog"
)

// HTMLTemplate -
type HTMLTemplate struct {
	tmpl *template.Template
}

// NewHTMLTemplate -
func NewHTMLTemplate(name, text string) (*HTMLTemplate, error) {
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return nil, err
	}
	return &HTMLTemplate{
		tmpl: tmpl,
	}, nil
}

// MustNewHTMLTemplate -
func MustNewHTMLTemplate(name, text string) *HTMLTemplate {
	t, err := NewHTMLTemplate(name, text)
	if err != nil {
		nlog.Fatal("fail to create html template: ", err)
	}
	return t
}

// Execute - applies the template to the specified data object
func (t *HTMLTemplate) Execute(data interface{}) (string, error) {
	var sb strings.Builder
	if err := t.tmpl.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// ExecuteTemplate - applies the template associated with t that has the given name to the specified data object
func (t *HTMLTemplate) ExecuteTemplate(name string, data interface{}) (string, error) {
	var sb strings.Builder
	if err := t.tmpl.ExecuteTemplate(&sb, name, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

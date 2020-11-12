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

package nutil

// Page -
type Page struct {
	PageNo   int32 `form:"pageNo" json:"pageNo" binding:"gt=0"`
	PageSize int32 `form:"pageSize" json:"pageSize" binding:"gt=0"`
	Total    int64 `json:"total"`
}

// NewPage -
func NewPage(pageNo int32, pageSize int32) *Page {
	return &Page{
		PageNo:   pageNo,
		PageSize: pageSize,
	}
}

// Offset -
func (p *Page) Offset() int {
	return int((p.PageNo - 1) * p.PageSize)
}

// Limit -
func (p *Page) Limit() int {
	return int(p.PageSize)
}

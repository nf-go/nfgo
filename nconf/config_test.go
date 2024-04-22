// Copyright 2024 The nfgo Authors. All Rights Reserved.
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

package nconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSensitiveURLPath(t *testing.T) {
	a := assert.New(t)
	c := &WebConfig{}
	isSensitive := c.IsSensitiveURLPath("/foo")
	a.False(isSensitive)
	c.SensitiveURLPaths = map[string]struct{}{
		"/login": {},
	}
	isSensitive = c.IsSensitiveURLPath("/login")
	a.True(isSensitive)
}

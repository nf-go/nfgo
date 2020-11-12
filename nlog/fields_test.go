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

package nlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFields(t *testing.T) {
	fs := NewFields()
	assert.Equal(t, 0, len(fs))

	fs = NewFields("k1")
	assert.Equal(t, 0, len(fs))

	fs = NewFields("k1", "v1")
	assert.Equal(t, 1, len(fs))
	assert.Equal(t, "v1", fs["k1"])

	fs = NewFields("k1", "v1", "k2")
	assert.Equal(t, 1, len(fs))
	assert.Equal(t, "v1", fs["k1"])

	fs = NewFields("k1", "v1", "k2", "v2")
	assert.Equal(t, 2, len(fs))
	assert.Equal(t, "v1", fs["k1"])
	assert.Equal(t, "v2", fs["k2"])

	fs = NewFields("k1", "", "k2", "")
	assert.Equal(t, 0, len(fs))
}

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
package ncrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256(t *testing.T) {
	plaint := "qwer1234JSQotTd2s3"
	hashed := Sha256(plaint)
	assert.Equal(t, "003b12dfe2c657c572f3496a63af8e2250f465e96aa0ebfea8428c82e0b10ba0", hashed)
	t.Log(len(hashed))
}

func TestRandomString(t *testing.T) {
	str := RandString(6)
	t.Log(str)
	assert.Equal(t, 6, len(str))
}

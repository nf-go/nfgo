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

package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	jwtSecret := "VoCra8#GEBAbRl*+vos9UF@??gi8Oy"
	token, err := NewToken(jwtSecret, "admin", time.Now().Add(time.Hour*24*365*20), map[string]interface{}{"userId": 1})
	t.Log(token)
	assert.NoError(t, err)
}

func TestParseToken(t *testing.T) {
	jwtSecret := "VoCra8#GEBAbRl*+vos9UF@??gi8Oy"
	token := "eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjIyMjY3MTIzMDkuODY5ODM5Nywic3ViIjoiYWRtaW4iLCJ1c2VySWQiOjF9.2M1tC3oQrC0Ym7K9qfVVcLS-eyCM4bacVOyuh6Jj1yg"
	claims, err := ParseToken(jwtSecret, token)
	assert.NoError(t, err)
	t.Log(claims.Subject)
	t.Log(claims.Expires)
	t.Log(claims.Set)
}

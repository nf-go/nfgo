// Copyright 2022 The nfgo Authors. All Rights Reserved.
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

func TestNewAESKey(t *testing.T) {
	a := assert.New(t)
	key, err := NewAESKey(256)
	a.Nil(err)
	a.Equal(32, len(key))
}

func TestEncryptAndDecryptString(t *testing.T) {
	a := assert.New(t)
	keyText, err := NewAESKeyString(256)
	a.Nil(err)
	text := "hello world"
	encryptedText, err := AESEncryptString(text, keyText)
	a.Nil(err)
	decryptedText, err := AESDecryptString(encryptedText, keyText)
	a.Nil(err)
	a.Equal(text, decryptedText)

}

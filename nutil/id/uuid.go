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

package id

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// UUID - represents a UUID.
type UUID [16]byte

// NewUUID - generates a new uuid.
func NewUUID() (UUID, error) {
	var uuid [16]byte

	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		return [16]byte{}, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return uuid, nil
}

// Hex - returns a hex representation of the UUID.
func (id UUID) Hex() string {
	return hex.EncodeToString(id[:])
}

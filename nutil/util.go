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

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"nfgo.ga/nfgo/nutil/id"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RedisKey -
type RedisKey string

// Key -
func (k RedisKey) Key(a ...interface{}) string {
	return fmt.Sprintf(string(k), a...)
}

// Sha256 -
func Sha256(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

// RandString -
func RandString(length int) string {
	b := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

// UUID -
func UUID() (string, error) {
	uuid, err := id.NewUUID()
	if err != nil {
		return "", err
	}
	return uuid.Hex(), err
}

// IsNil -
func IsNil(i interface{}) bool {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Interface, reflect.UnsafePointer:
		return v.IsNil()
	}
	return i == nil
}

// IsNotNil -
func IsNotNil(i interface{}) bool {
	return !IsNil(i)
}

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

package ntypes

import "testing"

func TestIsNil(t *testing.T) {
	var ch chan int
	if !IsNil(ch) {
		t.Fatal()
	}

	var f func() = nil
	if !IsNil(f) {
		t.Fatal()
	}

	var m map[string]int = nil
	if !IsNil(m) {
		t.Fatal()
	}

	var s []int = nil
	if !IsNil(s) {
		t.Fatal(s)
	}

	var obj interface{} = (*int)(nil)
	if !IsNil(obj) {
		t.Fatal()
	}

	var obj2 interface{} = nil
	if !IsNil(obj2) {
		t.Fatal()
	}

	if IsNil("") || IsNil(1) {
		t.Fatal()
	}

}

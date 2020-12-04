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

package nerrors

import "go.uber.org/multierr"

// Combine - Combine combines the passed errors into a single error.
// If zero arguments were passed or if all items are nil, a nil error is returned.
func Combine(errors ...error) error {
	return multierr.Combine(errors...)
}

// Append - appends the given errors together. Either value may be nil.
//
// This function is a specialization of Combine for the common case where
// there are only two errors.
func Append(left error, right error) error {
	return multierr.Append(left, right)
}

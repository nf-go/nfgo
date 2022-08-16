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

import "fmt"

var (
	// ErrInternal -
	ErrInternal = NewBizError(-1, "internal service error")
	// ErrUnauthorized -
	ErrUnauthorized = NewBizError(-2, "unauthorized")
	// ErrForbidden -
	ErrForbidden = NewBizError(-3, "forbidden")
)

// BizError -
type BizError interface {
	error
	Code() int
	Msg() string
	New(message string) BizError
}

// NewBizError -
func NewBizError(code int16, msg string) BizError {
	return &bizError{code: code, msg: msg}
}

// bizError -
type bizError struct {
	code int16
	msg  string
}

func (e *bizError) New(message string) BizError {
	return NewBizError(e.code, e.msg+": "+message)
}

func (e *bizError) Error() string {
	return fmt.Sprintf("biz error %d %s", e.code, e.msg)
}

// Code -
func (e *bizError) Code() int {
	return int(e.code)
}

// Msg -
func (e *bizError) Msg() string {
	return e.msg
}

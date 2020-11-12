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
	ErrInternal = NewBizError(-1, "内部服务错误")
	// ErrUnauthorized -
	ErrUnauthorized = NewBizError(-2, "未能通过认证")
	// ErrForbidden -
	ErrForbidden = NewBizError(-3, "无权访问")
)

// BizError -
type BizError interface {
	error
	Unwrap() error
	Code() int
	Msg() string
}

// NewBizError -
func NewBizError(code int16, msg string) BizError {
	return &bizError{code: code, msg: msg, err: nil}
}

// bizError - 业务错误
type bizError struct {
	code int16
	msg  string
	err  error
}

func (e *bizError) Error() string {
	return fmt.Sprintf("biz error %d %s", e.code, e.msg)
}

func (e *bizError) Unwrap() error {
	return e.err
}

// Code -
func (e *bizError) Code() int {
	return int(e.code)
}

// Msg -
func (e *bizError) Msg() string {
	return e.msg
}

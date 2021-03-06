// Copyright 2021 The nfgo Authors. All Rights Reserved.
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

import (
	"github.com/pkg/errors"
)

// New - returns an error with the supplied message.
// New also records the stack trace at the point it was called.
// In your application code, use nerrors.New or nerrros.Errorf to return errors.
func New(message string) error {
	return errors.New(message)
}

// Errorf - formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
// In your application code, use nerrors.New or nerrros.Errorf to return errors.
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
// If collaborating with other libraries, consider using nerrors.WithStack nerrors.Wrap or errors.Wrapf to store stack information. The same applies when working with standard libraries.
func WithStack(err error) error {
	return errors.WithStack(err)
}

// Wrap - returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
// If collaborating with other libraries, consider using nerrors.WithStack nerrors.Wrap or errors.Wrapf to store stack information. The same applies when working with standard libraries.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf - returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
// If collaborating with other libraries, consider using nerrors.WithStack nerrors.Wrap or errors.Wrapf to store stack information. The same applies when working with standard libraries.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

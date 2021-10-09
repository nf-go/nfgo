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

package ntime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartOfAndEndOf(t *testing.T) {
	a := assert.New(t)
	dt, err := time.Parse("2006-01-02 15:04:05", "2016-11-02 13:14:06")
	a.Nil(err)

	dt = StartOfDayUTC(dt)
	a.Equal("2016-11-02 00:00:00 +0000 UTC", dt.String())

	dt = EndOfDayUTC(dt)
	a.Equal("2016-11-02 23:59:59.999999999 +0000 UTC", dt.String())

	dt = StartOfDayLocal(dt)
	a.Contains(dt.String(), "2016-11-02 00:00:00")

	dt = EndOfDayLocal(dt)
	a.Contains(dt.String(), "2016-11-02 23:59:59.999999999")
}

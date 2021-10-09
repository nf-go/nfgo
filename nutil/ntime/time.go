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

import "time"

// StartOfDay -
func StartOfDay(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

// StartOfDayLocal -
func StartOfDayLocal(t time.Time) time.Time {
	return StartOfDay(t, time.Local)
}

// StartOfDayUTC -
func StartOfDayUTC(t time.Time) time.Time {
	return StartOfDay(t, time.UTC)
}

// EndOfDay -
func EndOfDay(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, loc)
}

// EndOfDayLocal -
func EndOfDayLocal(t time.Time) time.Time {
	return EndOfDay(t, time.Local)
}

// EndOfDayUTC -
func EndOfDayUTC(t time.Time) time.Time {
	return EndOfDay(t, time.UTC)
}

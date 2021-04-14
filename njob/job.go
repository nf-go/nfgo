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

package njob

import (
	"context"

	"nfgo.ga/nfgo/ncontext"
	"nfgo.ga/nfgo/nutil/ncrypto"
)

// Job -
type Job interface {
	Run(ctx context.Context) error
}

// FuncJob -
type FuncJob func(ctx context.Context) error

// Run -
func (f FuncJob) Run(ctx context.Context) error {
	return f(ctx)
}

// FuncJobs -
type FuncJobs map[string]FuncJob

// Jobs -
type Jobs map[string]Job

func newJobContext(jobName string) context.Context {
	mdc := ncontext.NewMDC()
	traceID, _ := ncrypto.UUID()
	mdc.SetTraceID(traceID)
	mdc.SetSubjectID(jobName)

	ctx := context.Background()
	return ncontext.WithMDC(ctx, mdc)
}

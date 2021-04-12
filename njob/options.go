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

type serverOptions struct {
	funcJobs         FuncJobs
	jobs             Jobs
	distributedMutex DistributedMutex
}

type ServerOption func(*serverOptions)

func FuncJobsOption(funcJobs FuncJobs) ServerOption {
	return func(opts *serverOptions) {
		opts.funcJobs = funcJobs
	}
}

func JobsOption(jobs Jobs) ServerOption {
	return func(opts *serverOptions) {
		opts.jobs = jobs
	}
}

func DistributedMutexOption(distributedMutex DistributedMutex) ServerOption {
	return func(opts *serverOptions) {
		opts.distributedMutex = distributedMutex
	}
}

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

package nconst

const (
	// HeaderTraceID -
	HeaderTraceID string = "X-Trace-ID"
	// HeaderRealIP -
	HeaderRealIP string = "X-Real-IP"
	// HeaderForwardedFor -
	HeaderForwardedFor string = "X-Forwarded-For"
	// HeaderToken -
	HeaderToken string = "X-Token"
	// HeaderSub -
	HeaderSub string = "X-Sub"
	// HeaderTs -
	HeaderTs string = "X-Ts"
	// HeaderSig - SHA256(signKey + X-Ts + X-Sub + X-Trace-ID)
	HeaderSig string = "X-Sig"
	// HeaderClientType -
	HeaderClientType string = "X-ClientType"
)

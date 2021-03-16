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

package ncontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextOthers(t *testing.T) {
	a := assert.New(t)
	m := NewMDC()
	v := m.Other("notexist")
	a.Nil(v)
}

func TestMDCCopy(t *testing.T) {
	m := NewMDC()
	m.SetAPIName("api")
	m.SetClientIP("ip")
	m.SetClientType("type")
	m.SetRPCName("rpc")
	m.SetSubjectID("s")
	m.SetTraceID("t")
	m.SetOther("k1", "v1")
	m.SetOther("k2", "v2")
	cm := m.Copy()

	a := assert.New(t)
	a.Equal("api", cm.APIName())
	a.Equal("ip", cm.ClientIP())
	a.Equal("type", cm.ClientType())
	a.Equal("rpc", cm.RPCName())
	a.Equal("s", cm.SubjectID())
	a.Equal("t", cm.TraceID())
	a.Equal("v1", cm.Other("k1"))
	a.Equal("v2", cm.Other("k2"))
	a.NotEqual(m, cm)
}

func TestBackground(t *testing.T) {
	m := NewMDC()
	m.SetAPIName("api")
	m.SetOther("k1", "v1")
	m.SetOther("k2", "v2")
	ctx1 := WithMDC(context.Background(), m)
	ctx2 := Background(ctx1)

	a := assert.New(t)
	a.NotEqual(ctx1, ctx2)
	cm, err := CurrentMDC(ctx2)
	a.Nil(err)
	a.Equal("api", cm.APIName())
	a.Equal("v1", cm.Other("k1"))
	a.Equal("v2", cm.Other("k2"))
	a.NotEqual(m, cm)
	a.Nil(ctx1.Value(0))
	a.NotNil(ctx1.Value(ctxKeyMDC))
}

func TestBackgroundNilMDC(t *testing.T) {

	ctx1 := context.Background()
	ctx2 := Background(ctx1)

	a := assert.New(t)
	a.Equal(ctx1, ctx2)
	cm, err := CurrentMDC(ctx2)
	a.Contains(err.Error(), "can't extract MDC from the context")
	a.Nil(cm)

}

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
	"errors"
	"sync"
)

type ctxKeyName int

const (
	// CtxKeyDb -
	CtxKeyDb ctxKeyName = iota
	// CtxKeyMDC -
	CtxKeyMDC
)

// ROMDC - Readonly Mapped Diagnostic Context
type ROMDC interface {
	ClientType() string
	TraceID() string
	SubjectID() string
	RPCName() string
	APIName() string
	ClientIP() string
	Other(key string) interface{}
}

// MDC - Mapped Diagnostic Context
type MDC interface {
	ROMDC
	SetClientType(clientType string)
	SetTraceID(traceID string)
	SetSubjectID(subjectID string)
	SetRPCName(rpcName string)
	SetAPIName(apiName string)
	SetClientIP(clinetIP string)
	SetOther(key string, value interface{})
}

// NewMDC -
func NewMDC() MDC {
	return &mdc{
		others: &sync.Map{},
	}
}

// WithMDC -
func WithMDC(ctx context.Context, m MDC) context.Context {
	return context.WithValue(ctx, CtxKeyMDC, m)
}

// CurrentMDC -
func CurrentMDC(ctx context.Context) (ROMDC, error) {
	v := ctx.Value(CtxKeyMDC)
	if mc, ok := v.(ROMDC); ok {
		return ROMDC(mc), nil
	}
	return nil, errors.New("can't extract MDC from the context")
}

type mdc struct {
	clientType string
	traceID    string
	subjectID  string
	rpcName    string
	apiName    string
	clinetIP   string
	others     *sync.Map
}

func (m *mdc) ClientType() string {
	return m.clientType
}

func (m *mdc) SetClientType(clientType string) {
	m.clientType = clientType
}

func (m *mdc) TraceID() string {
	return m.traceID
}

func (m *mdc) SetTraceID(traceID string) {
	m.traceID = traceID
}

func (m *mdc) SubjectID() string {
	return m.subjectID
}

func (m *mdc) SetSubjectID(subjectID string) {
	m.subjectID = subjectID
}

func (m *mdc) RPCName() string {
	return m.rpcName
}
func (m *mdc) SetRPCName(rpcName string) {
	m.rpcName = rpcName
}

func (m *mdc) APIName() string {
	return m.apiName
}

func (m *mdc) SetAPIName(apiName string) {
	m.apiName = apiName
}

func (m *mdc) ClientIP() string {
	return m.clinetIP
}

func (m *mdc) SetClientIP(clinetIP string) {
	m.clinetIP = clinetIP
}

func (m *mdc) Other(key string) interface{} {
	if v, ok := m.others.Load(key); ok {
		return v
	}
	return nil
}

func (m *mdc) SetOther(key string, value interface{}) {
	m.others.Store(key, value)
}

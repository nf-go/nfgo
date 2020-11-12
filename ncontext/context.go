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

// MDC - Mapped Diagnostic Context
type MDC interface {
	ClientType() string
	SetClientType(clientType string)

	TraceID() string
	SetTraceID(traceID string)

	SubjectID() string
	SetSubjectID(subjectID string)

	RPCName() string
	SetRPCName(rpcName string)

	APIName() string
	SetAPIName(apiName string)

	ClientIP() string
	SetClientIP(clinetIP string)

	Other(key string) interface{}
	SetOther(key string, value interface{})
}

// NewMDC -
func NewMDC() MDC {
	return &mdc{
		others: &sync.Map{},
	}
}

// BindMDCToContext -
func BindMDCToContext(ctx context.Context, m MDC) context.Context {
	return context.WithValue(ctx, CtxKeyMDC, m)
}

// CurrentMDC -
func CurrentMDC(ctx context.Context) (MDC, error) {
	v := ctx.Value(CtxKeyMDC)
	if mc, ok := v.(MDC); ok {
		return MDC(mc), nil
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

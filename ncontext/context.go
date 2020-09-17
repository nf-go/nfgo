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

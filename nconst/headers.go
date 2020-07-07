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

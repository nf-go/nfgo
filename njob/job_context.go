package njob

import (
	"context"

	"nfgo.ga/nfgo/ncontext"
	"nfgo.ga/nfgo/nutil"
)

// NewJobContext -
func NewJobContext(jobName string) context.Context {
	mdc := ncontext.NewMDC()
	traceID, _ := nutil.UUID()
	mdc.SetTraceID(traceID)
	mdc.SetSubjectID(jobName)

	ctx := context.Background()
	return ncontext.BindMDCToContext(ctx, mdc)
}

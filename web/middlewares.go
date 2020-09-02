package web

import (
	"bytes"
	"io"
	"io/ioutil"

	"nfgo.ga/nfgo/nconst"
	"nfgo.ga/nfgo/ncontext"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil"
)

// BindMDC - BindMDC MiddleWare
func BindMDC() HandlerFunc {
	return func(c *Context) {
		mdc := ncontext.NewMDC()

		traceID := c.GetHeader(nconst.HeaderTraceID)
		if traceID == "" {
			var err error
			if traceID, err = nutil.UUID(); err != nil {
				c.Fail(err)
				c.Abort()
				return
			}
		}
		mdc.SetTraceID(traceID)
		mdc.SetAPIName(c.Request.Method + " " + c.Request.URL.Path)
		mdc.SetClientIP(c.ClientIP())
		mdc.SetClientType(c.GetHeader(nconst.HeaderClientType))
		mdc.SetSubjectID(c.GetHeader(nconst.HeaderSub))

		ctx := ncontext.BindMDCToContext(c.Request.Context(), mdc)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Logging -
func Logging() HandlerFunc {
	return func(c *Context) {
		if c.IsMultipartReq() {
			nlog.Logger(c).WithField("req", c.Request.URL.RawQuery).Info()
		} else {
			var buf bytes.Buffer
			teeReader := io.TeeReader(c.Request.Body, &buf)
			body, _ := ioutil.ReadAll(teeReader)
			c.Request.Body = ioutil.NopCloser(&buf)
			nlog.Logger(c).WithField("req", c.Request.URL.RawQuery+" "+string(body)).Info()
		}

		c.Next()
	}
}

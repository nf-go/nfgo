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

package web

import (
	"bytes"
	"io"

	"github.com/nf-go/nfgo/ncontext"
	"github.com/nf-go/nfgo/nlog"
	"github.com/nf-go/nfgo/nutil/nconst"
	"github.com/nf-go/nfgo/nutil/ncrypto"
)

// BindMDC - BindMDC MiddleWare
func BindMDC() HandlerFunc {
	return func(c *Context) {
		mdc := ncontext.NewMDC()

		traceID := c.GetHeader(nconst.HeaderTraceID)
		if traceID == "" {
			var err error
			if traceID, err = ncrypto.UUID(); err != nil {
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

		ctx := ncontext.WithMDC(c.Request.Context(), mdc)
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
			body, _ := io.ReadAll(teeReader)
			c.Request.Body = io.NopCloser(&buf)
			nlog.Logger(c).WithField("req", c.Request.URL.RawQuery+" "+string(body)).Info()
		}

		c.Next()
	}
}

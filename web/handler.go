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
	"github.com/gin-gonic/gin"
	"github.com/nf-go/nfgo/nconf"
)

// HandlerFunc -
type HandlerFunc func(c *Context)

// WrapHandler -
func (h HandlerFunc) WrapHandler(conf *nconf.WebConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			Context:   c,
			webConfig: conf,
		}
		h(ctx)
	}
}

func toGinHandlers(conf *nconf.WebConfig, handlers ...HandlerFunc) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		ginHandlers = append(ginHandlers, h.WrapHandler(conf))
	}
	return ginHandlers
}

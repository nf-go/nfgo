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
	"nfgo.ga/nfgo/nmetrics"
)

type serverOptions struct {
	metricsServer nmetrics.Server
	middlewares   []HandlerFunc
}

func (opts *serverOptions) setMiddlewaresToEngine(engine *gin.Engine) {
	middleWares := []gin.HandlerFunc{gin.Recovery()}
	if opts.metricsServer != nil {
		middleWares = append(middleWares, opts.metricsServer.WebMetricsMiddleware())
	}
	middleWares = append(middleWares, BindMDC().WrapHandler(), Logging().WrapHandler())
	if len(opts.middlewares) > 0 {
		for _, m := range opts.middlewares {
			middleWares = append(middleWares, m.WrapHandler())
		}
	}
	engine.Use(middleWares...)
}

// ServerOption -
type ServerOption func(*serverOptions)

// MetricsServerOption -
func MetricsServerOption(s nmetrics.Server) ServerOption {
	return func(opts *serverOptions) {
		opts.metricsServer = s
	}
}

// MiddlewaresOption -
func MiddlewaresOption(middleware ...HandlerFunc) ServerOption {
	return func(opts *serverOptions) {
		opts.middlewares = middleware
	}
}

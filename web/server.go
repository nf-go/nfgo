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
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil/graceful"
)

// Server -
type Server interface {
	graceful.ShutdownServer

	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

type server struct {
	engine     *gin.Engine
	config     *nconf.Config
	httpServer *http.Server
}

func (s *server) Serve() error {
	nlog.Info("the web server is started and serving on ", s.httpServer.Addr)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		nlog.Error("the web server is stoped  with error ", err)
		return err
	}
	return nil
}

func (s *server) MustServe() {
	if err := s.Serve(); err != nil {
		nlog.Fatal("fail to start http server:", err)
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *server) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	ginHandlers := toGinHandlers(handlers...)
	ginGroup := s.engine.Group(relativePath, ginHandlers...)
	return &routerGroup{ginGroup: ginGroup}
}

func (s *server) configSwagger() error {
	webConf := s.config.Web
	swaggerConf := webConf.Swagger
	if swaggerConf != nil && swaggerConf.Enabled {
		swagURL := swaggerConf.URL
		if swagURL == "" {
			if webConf.Host == "0.0.0.0" {
				swagURL = fmt.Sprintf("http://127.0.0.1:%d/apidoc/doc.json", webConf.Port)
			} else {
				swagURL = fmt.Sprintf("http://%s:%d/apidoc/doc.json", webConf.Host, webConf.Port)
			}
		}

		u, err := url.Parse(swagURL)
		if err != nil {
			return fmt.Errorf("fail to parse swagger url: %w", err)
		}

		relativePath := strings.ReplaceAll(u.Path, "/doc.json", "")

		s.engine.GET(relativePath+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(swagURL)))
	}
	return nil
}

// NewServer -
func NewServer(config *nconf.Config, opt ...ServerOption) (Server, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if config.App.IsProfileLocal() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	webConfig := config.Web
	if webConfig == nil {
		return nil, errors.New("web config is not initialized in the config")
	}

	// gin engine
	engine := gin.New()
	engine.MaxMultipartMemory = webConfig.MaxMultipartMemory
	opts := &serverOptions{}
	for _, o := range opt {
		o(opts)
	}
	opts.setMiddlewaresToEngine(engine)

	// http server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", webConfig.Host, webConfig.Port),
		Handler: engine,
	}

	s := &server{
		engine:     engine,
		config:     config,
		httpServer: httpServer,
	}

	// config swagger
	if err := s.configSwagger(); err != nil {
		return nil, err
	}

	return s, nil
}

// MustNewServer -
func MustNewServer(config *nconf.Config, opt ...ServerOption) Server {
	server, err := NewServer(config, opt...)
	if err != nil {
		nlog.Fatal("fail to init http server: ", err)
	}
	return server
}

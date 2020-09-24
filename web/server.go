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
	"nfgo.ga/nfgo/ngrace"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nmetrics"
)

// Server -
type Server interface {
	ngrace.Server
	RegisterOnShutdown(f func())
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

// ServerOption -
type ServerOption struct {
	MetricsServer nmetrics.Server
	Middlewares   []HandlerFunc
}

func (o *ServerOption) setMiddlewaresToEngine(engine *gin.Engine) {
	middleWares := []gin.HandlerFunc{gin.Recovery()}
	if o.MetricsServer != nil {
		middleWares = append(middleWares, o.MetricsServer.WebMetricsMiddleware())
	}
	middleWares = append(middleWares, BindMDC().WrapHandler(), Logging().WrapHandler())
	if len(o.Middlewares) > 0 {
		for _, m := range o.Middlewares {
			middleWares = append(middleWares, m.WrapHandler())
		}
	}
	engine.Use(middleWares...)
}

type server struct {
	engine *gin.Engine
	config *nconf.Config
	// host       string
	// port       int32
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

func (s *server) RegisterOnShutdown(f func()) {
	s.httpServer.RegisterOnShutdown(f)
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
func NewServer(config *nconf.Config, option *ServerOption) (Server, error) {
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
	if option == nil {
		option = &ServerOption{}
	}
	option.setMiddlewaresToEngine(engine)

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
func MustNewServer(config *nconf.Config, option *ServerOption) Server {
	server, err := NewServer(config, option)
	if err != nil {
		nlog.Fatal("fail to init http server: ", err)
	}
	return server
}

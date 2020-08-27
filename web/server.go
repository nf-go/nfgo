package web

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
)

// Server -
type Server interface {
	Run() error
	MustRun()
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

type server struct {
	engine *gin.Engine
	config *nconf.Config
	host   string
	port   int32
}

func (s *server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	nlog.Info("the web server is started and serving on ", addr)
	err := s.engine.Run(addr)
	if err != nil {
		nlog.Error("the web server is stoped  with error ", err)
	}
	return err
}

func (s *server) MustRun() {
	if err := s.Run(); err != nil {
		nlog.Fatal("fail to start http server:", err)
	}
}

func (s *server) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	ginHandlers := toGinHandlers(handlers...)
	ginGroup := s.engine.Group(relativePath, ginHandlers...)
	return &routerGroup{ginGroup: ginGroup}
}

// NewServer -
func NewServer(config *nconf.Config, middleware ...HandlerFunc) (Server, error) {
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

	engine := gin.New()

	engine.MaxMultipartMemory = webConfig.MaxMultipartMemory

	if len(middleware) == 0 {
		engine.Use(
			gin.Recovery(),
			BindMDC(config).WrapHandler(),
			Logging().WrapHandler(),
		)
	} else {
		engine.Use(gin.Recovery())
		for _, m := range middleware {
			engine.Use(m.WrapHandler())
		}
	}

	s := &server{
		engine: engine,
		config: config,
		host:   webConfig.Host,
		port:   webConfig.Port,
	}

	// config swagger
	swaggerConf := webConfig.Swagger
	if swaggerConf != nil && swaggerConf.Enabled {

		swagURL := swaggerConf.URL
		if swagURL == "" {
			if s.host == "0.0.0.0" {
				swagURL = fmt.Sprintf("http://127.0.0.1:%d/apidoc/doc.json", s.port)
			} else {
				swagURL = fmt.Sprintf("http://%s:%d/apidoc/doc.json", s.host, s.port)
			}
		}

		u, err := url.Parse(swagURL)
		if err != nil {
			return nil, fmt.Errorf("fail to parse swagger url: %w", err)
		}

		relativePath := strings.ReplaceAll(u.Path, "/doc.json", "")

		engine.GET(relativePath+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(swagURL)))
	}

	return s, nil
}

// MustNewServer -
func MustNewServer(config *nconf.Config, middleware ...HandlerFunc) Server {
	server, err := NewServer(config, middleware...)
	if err != nil {
		nlog.Fatal("fail to init http server: ", err)
	}
	return server
}

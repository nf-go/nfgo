package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RouterGroup -
type RouterGroup interface {
	Routes
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

type routerGroup struct {
	ginGroup *gin.RouterGroup
}

func (g *routerGroup) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	ginHandlers := toGinHandlers(handlers...)
	ginGroup := g.ginGroup.Group(relativePath, ginHandlers...)
	return &routerGroup{ginGroup: ginGroup}
}

func (g *routerGroup) Use(handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.Use(ginHandlers...)
}

func (g *routerGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.Handle(httpMethod, relativePath, ginHandlers...)
}

func (g *routerGroup) Any(relativePath string, handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.Any(relativePath, ginHandlers...)
}

func (g *routerGroup) GET(relativePath string, handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.GET(relativePath, ginHandlers...)
}

func (g *routerGroup) POST(relativePath string, handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.POST(relativePath, ginHandlers...)
}

func (g *routerGroup) DELETE(relativePath string, handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.DELETE(relativePath, ginHandlers...)
}

func (g *routerGroup) PATCH(relativePath string, handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.PATCH(relativePath, ginHandlers...)
}

func (g *routerGroup) PUT(relativePath string, handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.PUT(relativePath, ginHandlers...)
}

func (g *routerGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.OPTIONS(relativePath, ginHandlers...)

}

func (g *routerGroup) HEAD(relativePath string, handlers ...HandlerFunc) {
	ginHandlers := toGinHandlers(handlers...)
	g.ginGroup.HEAD(relativePath, ginHandlers...)

}

func (g *routerGroup) StaticFile(relativePath string, filepath string) {
	g.ginGroup.StaticFile(relativePath, filepath)
}

func (g *routerGroup) Static(relativePath string, root string) {
	g.ginGroup.Static(relativePath, root)

}

func (g *routerGroup) StaticFS(relativePath string, fs http.FileSystem) {
	g.ginGroup.StaticFS(relativePath, fs)
}

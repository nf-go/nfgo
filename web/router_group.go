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

// RouterRegistrar
type RouterRegistrar interface {
	RegisterRoutes(rg RouterGroup)
}

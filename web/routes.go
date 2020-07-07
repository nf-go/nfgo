package web

import "net/http"

// Routes defines all router handle interface.
type Routes interface {
	Use(handlers ...HandlerFunc)

	Handle(httpMethod, relativePath string, handlers ...HandlerFunc)

	Any(relativePath string, handlers ...HandlerFunc)

	GET(relativePath string, handlers ...HandlerFunc)

	POST(relativePath string, handlers ...HandlerFunc)

	DELETE(relativePath string, handlers ...HandlerFunc)

	PATCH(relativePath string, handlers ...HandlerFunc)

	PUT(relativePath string, handlers ...HandlerFunc)

	OPTIONS(relativePath string, handlers ...HandlerFunc)

	HEAD(relativePath string, handlers ...HandlerFunc)

	StaticFile(relativePath string, filepath string)

	Static(relativePath string, root string)

	StaticFS(relativePath string, fs http.FileSystem)
}

package web

import "github.com/gin-gonic/gin"

// HandlerFunc -
type HandlerFunc func(c *Context)

// WrapHandler -
func (h HandlerFunc) WrapHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			c,
		}
		h(ctx)
	}
}

func toGinHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		ginHandlers = append(ginHandlers, h.WrapHandler())
	}
	return ginHandlers
}

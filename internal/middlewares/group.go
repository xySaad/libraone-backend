package middlewares

import "github.com/gin-gonic/gin"

type MiddlewareFunc[D any] func(c *gin.Context) *D

type Middleware[D any] struct {
	gin.RouterGroup
	middleware MiddlewareFunc[D]
}

func Group[D any](routerGroup gin.RouterGroup, relativePath string, middleware MiddlewareFunc[D]) Middleware[D] {
	return Middleware[D]{RouterGroup: routerGroup, middleware: middleware}
}

type HandlerFunc[D any] func(ctx *gin.Context, dependency *D)

func (m *Middleware[D]) GET(relativePath string, handler HandlerFunc[D]) gin.IRoutes {
	return m.RouterGroup.GET(relativePath, func(ctx *gin.Context) {
		result := m.middleware(ctx)
		if result != nil {
			handler(ctx, result)
		}
	})
}

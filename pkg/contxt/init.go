package contxt

import (
	"context"

	"github.com/gin-gonic/gin"
)

// SetupContext is a middleware to embbed this Context type into gin.Context
func SetupContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), tCtxKey, initContext(c))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// initContext creates a basic context from gin.Context
func initContext(c *gin.Context) *Context {
	var requestID, prefix string
	if reqID, ok := c.Get("RequestID"); ok {
		requestID, _ = reqID.(string)
	}
	if px, ok := c.Get("LogPrefix"); ok {
		prefix, _ = px.(string)
	}

	ctx := &Context{
		ginCtx:  c,
		prefix:  prefix,
		reqInfo: RequestInfo{RequestID: requestID},
	}
	if c.Request != nil {
		ctx.reqInfo = RequestInfo{
			RequestID: ctx.reqInfo.RequestID,
			ClientIP:  c.ClientIP(),
			Host:      c.Request.Host,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Referer:   c.Request.Referer(),
			UserAgent: c.Request.UserAgent(),
		}
		ctx.parentCtx = c.Request.Context()
	}
	return ctx
}

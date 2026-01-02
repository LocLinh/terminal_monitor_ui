package contxt

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ctxKey int

const (
	tCtxKey ctxKey = iota
)

type Context struct {
	ginCtx    *gin.Context
	prefix    string
	reqInfo   RequestInfo
	parentCtx context.Context
}

// RequestInfo for Context
type RequestInfo struct {
	RequestID string
	ClientIP  string
	Host      string
	Method    string
	Path      string
	Referer   string
	UserAgent string
}

// IContext is an interface for Context
type IContext interface {
	//set
	Set(key string, val interface{})

	//get
	GetLogger() IContextLogger
	GetLoggerWithPrefix(prefix string) IContextLogger
	Get(key string) (val interface{}, ok bool)
	GetString(key string) (val string, ok bool)
	GetInt64(key string) (val int64, ok bool)
	GetRequestID() string
	GetRequestInfo() RequestInfo

	//request
	Request() *http.Request
	RequestStart()
	RequestFinished()
	RequestFinishedWithResponseBody(responseBody string)

	//context
	Context() context.Context
}

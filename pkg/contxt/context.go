package contxt

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"runtime"
	"time"

	"terminal_monitor_ui/pkg/wrapper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewContext is used to get Context out from the regular context.Context
func NewContext(c context.Context) IContext {
	ginWrapperCtx, ok := c.(*wrapper.Context)
	if ok {
		return NewContext(ginWrapperCtx.Context)
	}

	ginCtx, ok := c.(*gin.Context)
	if ok {
		c = ginCtx.Request.Context()
	}
	value := c.Value(tCtxKey)
	if value == nil {
		return newContextWithRandomRequestID()
	}

	ctx, ok := value.(*Context)
	if !ok {
		return newContextWithRandomRequestID()
	}

	ctx.parentCtx = c
	return ctx
}

// NewContextFromRequestID...
func NewContextFromRequestID(reqID string) IContext {
	return &Context{reqInfo: RequestInfo{RequestID: reqID}}
}

// GetLogger ...
func (c *Context) GetLogger() IContextLogger {
	return &ContextLogger{reqInfo: c.reqInfo, context: c.Context()}
}

// GetLoggerWithPrefix ...
func (c *Context) GetLoggerWithPrefix(prefix string) IContextLogger {
	return &ContextLogger{
		prefix:  prefix,
		reqInfo: c.reqInfo,
		context: c.Context(),
	}
}

// Context returns an empty context ONLY for logging purposes
func (c *Context) Context() context.Context {
	if c.parentCtx != nil {
		return c.parentCtx
	}
	return context.WithValue(context.Background(), tCtxKey, c)
}

// Set ...
func (c *Context) Set(key string, value interface{}) {
	c.ginCtx.Set(key, value)
}

// Get ...
func (c *Context) Get(key string) (interface{}, bool) {
	return c.ginCtx.Get(key)
}

// Request return the current request
func (c *Context) Request() *http.Request {
	return c.ginCtx.Request
}

// GetString ...
func (c *Context) GetString(key string) (string, bool) {
	v, ok := c.Get(key)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return s, true
}

// GetInt64 ...
func (c *Context) GetInt64(key string) (int64, bool) {
	v, ok := c.Get(key)
	if !ok {
		return 0, false
	}
	i, ok := v.(int64)
	if !ok {
		return 0, false
	}
	return i, true
}

// GetRequestID ...
func (c *Context) GetRequestID() string {
	return c.reqInfo.RequestID
}

// GetRequestInfo ...
func (c *Context) GetRequestInfo() RequestInfo {
	return c.reqInfo
}

// RequestStart mark the start of the current request and
// ship request data to Stackdriver
func (c *Context) RequestStart() {
	c.Set("RequestTime", time.Now())
	c.GetLogger().SetPrefix("Middleware").Infof("Request started")
}

// RequestFinished output the marker for current request and
// ship response data to Stackdriver
func (c *Context) RequestFinished() {
	c.RequestFinishedWithResponseBody("")
}

// RequestFinishedWithResponseBody output the marker for current request and
// ship response data to Stackdriver
func (c *Context) RequestFinishedWithResponseBody(responseBody string) {
	statusCode := c.ginCtx.Writer.Status()
	s, _ := c.Get("RequestTime")
	reqTime, _ := s.(time.Time)
	latency := int(math.Ceil(float64(time.Since(reqTime).Nanoseconds()) / 1e6))

	c.prefix = "Middleware"
	responseLength := fmt.Sprintf("%dB", c.ginCtx.Writer.Size())
	fmtLatency := fmt.Sprintf("%dms", latency)

	fields := c.basicFields()
	exfields := []zapcore.Field{
		zap.Int("statusCode", statusCode),
		zap.String("latency", fmtLatency),
		zap.String("responseLength", responseLength),
	}

	fields = append(fields, exfields...)
	if len(responseBody) > 0 {
		fields = append(fields, zap.String("response_body", responseBody))
	}

	// assign log level
	if statusCode > 499 {
		zap.L().Error("Request finished", fields...)
	} else if statusCode > 399 {
		zap.L().Warn("Request finished", fields...)
	} else {
		zap.L().Info("Request finished", fields...)
	}
}

// support func to build up basic fields
func (c *Context) basicFields() []zapcore.Field {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	// frame of getCurrentFrame function
	frame, _ := frames.Next() // nolint
	// frame of caller of this function
	frame, _ = frames.Next() // nolint
	logLine := fmt.Sprintf("%s:%d", frame.File, frame.Line)

	fields := []zapcore.Field{
		zap.String("prefix", c.prefix),
		zap.String("request_id", c.reqInfo.RequestID),
		zap.String("client_ip", c.reqInfo.ClientIP),
		zap.String("host", c.reqInfo.Host),
		zap.String("method", c.reqInfo.Method),
		zap.String("path", c.reqInfo.Path),
		zap.String("referer", c.reqInfo.Referer),
		zap.String("user_agent", c.reqInfo.UserAgent),
		zap.String("log_line", logLine),
	}

	return fields
}

func GetWrapper(ctx context.Context) (*Context, error) {
	value := ctx.Value(tCtxKey)
	if value == nil {
		return nil, errors.New("could not get wrapper.Context from context")
	}

	wrapperCtx, ok := value.(*Context)
	if !ok {
		return nil, errors.New("could not get wrapper.Context from context")
	}

	return wrapperCtx, nil
}

func Background() context.Context {
	ctx := context.Background()
	return ctxzap.ToContext(ctx, zap.L())
}

func newContextWithRandomRequestID() *Context {
	uuid, _ := uuid.NewUUID()
	return &Context{reqInfo: RequestInfo{RequestID: uuid.String()}}
}

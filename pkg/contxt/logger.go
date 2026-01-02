package contxt

import (
	"context"
	"fmt"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewContextLoggerFromReqInfo ...
func NewContextLoggerFromReqInfo(reqInfo RequestInfo) *ContextLogger {
	ctx := NewContext(context.Background()).(*Context)
	ctx.reqInfo = reqInfo
	return &ContextLogger{reqInfo: reqInfo, context: ctx.Context()}
}

func NewLoggerWithPrefix(ctx context.Context, prefix string) IContextLogger {
	wctx, _ := GetWrapper(ctx)
	return &ContextLogger{prefix: prefix, context: ctx, reqInfo: wctx.GetRequestInfo()}
}

func (l *ContextLogger) Context() context.Context {
	return l.context
}

// Infof log formatted info with requestID
func (l *ContextLogger) Infof(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	if len(str) > 2048 {
		str = str[0:2047]
	}
	zap.L().Info(str, l.basicFields()...)
}

// Warnf log formatted warn with requestID
func (l *ContextLogger) Warnf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	zap.L().Warn(str, l.basicFields()...)
}

// Errorf log formatted error with requestID
func (l *ContextLogger) Errorf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	zap.L().Warn(str, l.basicFields()...)
}

// SetPrefix ...
func (l *ContextLogger) SetPrefix(prefix string) IContextLogger {
	l.prefix = prefix
	return l
}

// SetField ...
func (l *ContextLogger) SetField(k string, v interface{}) IContextLogger {
	l.fields = append(l.fields, zap.String(k, fmt.Sprintf("%+v", v)))
	return l
}

// support func to build up basic fields
func (l *ContextLogger) basicFields() []zapcore.Field {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	// frame of getCurrentFrame function
	frame, _ := frames.Next() // nolint
	// frame of caller of this function
	frame, _ = frames.Next() // nolint
	logLine := fmt.Sprintf("%s:%d", frame.File, frame.Line)

	fields := []zapcore.Field{
		zap.String("prefix", l.prefix),
		zap.String("request_id", l.reqInfo.RequestID),
		zap.String("client_ip", l.reqInfo.ClientIP),
		zap.String("host", l.reqInfo.Host),
		zap.String("method", l.reqInfo.Method),
		zap.String("path", l.reqInfo.Path),
		zap.String("referer", l.reqInfo.Referer),
		zap.String("user_agent", l.reqInfo.UserAgent),
		zap.String("log_line", logLine),
	}

	return append(l.fields, fields...)
}

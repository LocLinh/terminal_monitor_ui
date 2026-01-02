package contxt

import (
	"context"

	"go.uber.org/zap/zapcore"
)

type IContextLogger interface {
	SetField(string, interface{}) IContextLogger
	SetPrefix(prefix string) IContextLogger
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Context() context.Context
}

type ContextLogger struct {
	prefix  string
	reqInfo RequestInfo
	fields  []zapcore.Field
	context context.Context
}

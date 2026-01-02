package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ChannelWriter struct {
	ch chan string
}

func NewChannelWriter(ch chan string) *ChannelWriter {
	return &ChannelWriter{ch: ch}
}

func (w *ChannelWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSuffix(string(p), "\n")
	select {
	case w.ch <- msg:
	default:
		// Drop message if channel is full
	}
	return len(p), nil
}

func (w *ChannelWriter) Sync() error {
	return nil
}

func NewLogger(ch chan string) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	writer := NewChannelWriter(ch)
	core := zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel)

	return zap.New(core)
}

package aiyaml

import (
	"context"
)

// Logger 日志接口
type Logger interface {
	WithContext(ctx context.Context) Logger
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

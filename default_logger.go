package aiyaml

import (
	"context"
	"log"
)

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	ctx context.Context
}

// NewDefaultLogger 创建默认日志器
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

// WithContext 添加上下文
func (dl *DefaultLogger) WithContext(ctx context.Context) Logger {
	return &DefaultLogger{ctx: ctx}
}

// WithField 添加字段
func (dl *DefaultLogger) WithField(key string, value interface{}) Logger {
	return dl
}

// WithError 添加错误
func (dl *DefaultLogger) WithError(err error) Logger {
	return dl
}

// Info 信息日志
func (dl *DefaultLogger) Info(args ...interface{}) {
	log.Println(args...)
}

// Infof 格式化信息日志
func (dl *DefaultLogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Error 错误日志
func (dl *DefaultLogger) Error(args ...interface{}) {
	log.Println(args...)
}

// Errorf 格式化错误日志
func (dl *DefaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

package vela

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log interface {
	LoggerLevel() zapcore.Level
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Panic(...interface{})
	Fatal(...interface{})
	Trace(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
	Tracef(string, ...interface{})
	Replace(*zap.Logger)
}

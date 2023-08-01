package logger

import (
	"errors"
	"github.com/vela-ssoc/vela-kit/lua"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	state             *zapState
	errRequiredOutput = errors.New("至少选择一种日志输出方式")
)

func init() {
	//初始化sate
	state = newZapState(defaultConfig())

	//Errorf("init logger %v succeed" , state.cfg )
}

type zapCallback func() error

type zapState struct {
	lua.SuperVelaData

	cfg    *config
	stop   zapCallback
	rotate zapCallback

	sugar *zap.SugaredLogger
}

func (z *zapState) newSugar() {
	sugar, stopFn, rotate := newSugar(z.cfg)
	z.sugar = sugar
	z.stop = stopFn
	z.rotate = rotate
}

func (z *zapState) Error(args ...interface{}) {
	z.sugar.Error(args...)
}

func (z *zapState) Debug(args ...interface{}) {
	z.sugar.Debug(args...)
}

func (z *zapState) Info(args ...interface{}) {
	z.sugar.Info(args...)
}

func (z *zapState) Warn(args ...interface{}) {
	z.sugar.Warn(args...)
}

func (z *zapState) Trace(args ...interface{}) {
	z.sugar.Debug(args...)
}

func (z *zapState) Panic(args ...interface{}) {
	z.sugar.Panic(args...)
}

func (z *zapState) Fatal(args ...interface{}) {
	z.sugar.Panic(args...)
}

func (z *zapState) Errorf(format string, args ...interface{}) {
	z.sugar.Errorf(format, args...)
}

func (z *zapState) Debugf(format string, args ...interface{}) {
	z.sugar.Debugf(format, args...)
}

func (z *zapState) Infof(format string, args ...interface{}) {
	z.sugar.Infof(format, args...)
}

func (z *zapState) Warnf(format string, args ...interface{}) {
	z.sugar.Warnf(format, args...)
}

func (z *zapState) Panicf(format string, args ...interface{}) {
	z.sugar.Panicf(format, args...)
}

func (z *zapState) Fatalf(format string, args ...interface{}) {
	z.sugar.Panicf(format, args...)
}

func (z *zapState) Tracef(format string, args ...interface{}) {
	z.sugar.Debugf(format, args...)
}

func (z *zapState) Replace(*zap.Logger) {
}

func (z *zapState) LoggerLevel() zapcore.Level {
	return z.sugar.Level()
}

func newZapState(cfg *config) *zapState {
	obj := &zapState{
		cfg:  cfg,
		stop: func() error { return nil },
	}
	obj.newSugar()

	return obj
}

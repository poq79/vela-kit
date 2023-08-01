package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func encoder(format string, color bool) zapcore.Encoder {

	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	switch format {
	case FormatJson:
		c.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewJSONEncoder(c)
	default:
		if color {
			c.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		return zapcore.NewConsoleEncoder(c)
	}
}

func newSugar(c *config) (*zap.SugaredLogger, zapCallback, zapCallback) {

	//初始化stop函数
	stopFn := func() error {
		return nil
	}

	rotate := func() error {
		return nil
	}

	//必要函数和等级要求
	encode := encoder(c.Format, c.Color)
	fn := zap.LevelEnablerFunc(func(v zapcore.Level) bool {
		return v >= c.Level
	})

	var core zapcore.Core

	//输出到文件
	if c.Filename != "" {
		w := &lumberjack.Logger{
			Filename:   c.Filename,
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			Compress:   c.Compress,
		}
		core = zapcore.NewCore(encode, zapcore.AddSync(w), fn)
		stopFn = w.Close
		rotate = w.Rotate

	}

	//输出到前台
	if c.Console {
		sync := zapcore.AddSync(os.Stderr)
		if core == nil {
			core = zapcore.NewCore(encode, sync, fn)
		} else {
			core = zapcore.NewTee(core, zapcore.NewCore(encode, sync, fn))
		}
	}

	return zap.New(core, zap.AddCallerSkip(c.Skip), zap.WithCaller(c.Caller)).Sugar(), stopFn, rotate
}

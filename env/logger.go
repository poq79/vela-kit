package env

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (env *Environment) Debug(args ...interface{}) {
	env.log.Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func (env *Environment) Info(args ...interface{}) {
	env.log.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (env *Environment) Warn(args ...interface{}) {
	env.log.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (env *Environment) Error(args ...interface{}) {
	env.log.Error(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (env *Environment) Panic(args ...interface{}) {
	env.log.Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.CatchE.
func (env *Environment) Fatal(args ...interface{}) {
	//env.log.Fatal(args...)
	env.log.Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.CatchE.
func (env *Environment) Trace(args ...interface{}) {
	//env.log.Fatal(args...)
	env.log.Trace(args...)
}

func (env *Environment) Debugf(template string, args ...interface{}) {
	env.log.Debugf(template, args...)
}

func (env *Environment) Infof(template string, args ...interface{}) {
	env.log.Infof(template, args...)
}

func (env *Environment) Tracef(template string, args ...interface{}) {
	env.log.Tracef(template, args...)
}

func (env *Environment) Warnf(template string, args ...interface{}) {
	env.log.Warnf(template, args...)
}

func (env *Environment) Errorf(template string, args ...interface{}) {
	env.log.Errorf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (env *Environment) Panicf(template string, args ...interface{}) {
	env.log.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.CatchE.
func (env *Environment) Fatalf(template string, args ...interface{}) {
	//env.log.Fatalf(template, args...)
	env.log.Fatalf(template, args...)
}

func (env *Environment) Printf(msg string, keysAndValues ...interface{}) {
	env.log.Errorf(msg, keysAndValues...)
}

func (env *Environment) Replace(z *zap.Logger) {
	env.log.Replace(z)
}

func (env *Environment) LoggerLevel() zapcore.Level {
	return env.log.LoggerLevel()
}

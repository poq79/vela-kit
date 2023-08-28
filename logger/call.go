package logger

func Error(args ...interface{}) {
	xEnv.Error(args...)
}

func Debug(args ...interface{}) {
	xEnv.Debug(args...)
}

func Info(args ...interface{}) {
	xEnv.Info(args...)
}

func Warn(args ...interface{}) {
	xEnv.Warn(args...)
}

func Trace(args ...interface{}) {
	xEnv.Debug(args...)
}

func Panic(args ...interface{}) {
	xEnv.Panic(args...)
}

func Fatal(args ...interface{}) {
	xEnv.Panic(args...)
}

func Errorf(format string, args ...interface{}) {
	xEnv.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	xEnv.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	xEnv.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	xEnv.Warnf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	xEnv.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	xEnv.Panicf(format, args...)
}

func Tracef(format string, args ...interface{}) {
	xEnv.Debugf(format, args...)
}

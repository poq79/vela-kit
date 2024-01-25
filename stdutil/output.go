package stdutil

import (
	"fmt"
	"runtime"
	"time"
)

type Output struct {
	opt *Option
}

func New(args ...func(*Option)) *Output {
	opt := NewOption()
	for _, fn := range args {
		fn(opt)
	}

	return &Output{opt: opt}
}

func (out *Output) standard(format string, level string) string {

	var f string
	if out.opt.debug {
		_, file, line, ok := runtime.Caller(out.opt.skip)
		if ok {
			fn := fmt.Sprintf("[%s:%d]", file, line)
			f = time.Now().Format("2006-01-02 15:04:05") + " " + fn + " [" + level + "] " + format
		}
		return f + "\n"
	}

	f = time.Now().Format("2006-01-02 15:04:05") + " [" + level + "] " + format
	f = f + "\n"
	return f
}

func (out *Output) Write(data []byte) (int, error) {
	return out.opt.WriteCloser.Write(data)
}

func (out *Output) ERR(format string, v ...interface{}) {
	out.opt.WriteFn(out.standard(format, "ERROR"), v...)
}

func (out *Output) Info(format string, v ...interface{}) {
	out.opt.WriteFn(out.standard(format, "INFO"), v...)
}

func (out *Output) Warn(format string, v ...interface{}) {
	out.opt.WriteFn(out.standard(format, "WARN"), v...)
}

func (out *Output) Debug(format string, v ...interface{}) {
	out.opt.WriteFn(out.standard(format, "DEBUG"), v...)
}

func (out *Output) Fatal(format string, v ...interface{}) {
	out.opt.WriteFn(out.standard(format, "FATAL"), v...)
}

func (out *Output) Close() error {
	if out.opt.WriteCloser != nil {
		return out.opt.WriteCloser.Close()
	}
	return nil
}

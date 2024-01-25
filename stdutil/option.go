package stdutil

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/strutil"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

/*
	out := stdutil.New(stdutil.Print())
	out.Error()
	out.Info()
*/

type Option struct {
	debug       bool
	skip        int
	WriteFn     func(string, ...interface{})
	WriteCloser io.WriteCloser
}

func NewOption() *Option {
	opt := &Option{skip: 2, debug: false}
	Print()(opt)
	return opt
}

func Skip(n int) func(*Option) {
	return func(opt *Option) {
		opt.skip = n
	}
}

func Debug(flag bool) func(*Option) {
	return func(opt *Option) {
		opt.debug = flag
	}
}

func Print() func(*Option) {
	return func(opt *Option) {
		opt.WriteFn = func(format string, val ...interface{}) {
			fmt.Printf(format, val...)
		}
	}
}

func Console() func(*Option) {
	return func(opt *Option) {
		exe, _ := os.Executable()
		w := &lumberjack.Logger{
			Filename:   filepath.Join(filepath.Dir(exe), "console.log"),
			MaxSize:    50, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}

		opt.WriteFn = func(format string, val ...interface{}) {
			//header := time.Now().Format("2006-01-02 15:04:05") + " ERROR " + format
			//header = header + "\n"
			//if w == nil {
			//	fmt.Printf(header, val...)
			//	return
			//}
			_, _ = w.Write(strutil.S2B(fmt.Sprintf(format, val...)))
		}
		opt.WriteCloser = w
	}
}

func Upgrade() func(*Option) {
	return func(opt *Option) {
		exe, _ := os.Executable()
		w := &lumberjack.Logger{
			Filename:   filepath.Join(filepath.Dir(exe), "upgrade.log"),
			MaxSize:    50, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}

		opt.WriteFn = func(format string, val ...interface{}) {
			_, _ = w.Write(strutil.S2B(fmt.Sprintf(format, val...)))
		}
		opt.WriteCloser = w
	}
}

func Daemon() func(*Option) {
	return func(opt *Option) {
		exe, _ := os.Executable()
		w := &lumberjack.Logger{
			Filename:   filepath.Join(filepath.Dir(exe), "daemon.log"),
			MaxSize:    50, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}

		opt.WriteFn = func(format string, val ...interface{}) {
			_, _ = w.Write(strutil.S2B(fmt.Sprintf(format, val...)))
		}
		opt.WriteCloser = w
	}
}

package auxlib

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

func NewOutput(file string) io.WriteCloser {
	exe, _ := os.Executable()
	w := &lumberjack.Logger{
		Filename:   filepath.Join(filepath.Dir(exe), file),
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     7,    //days
		Compress:   true, // disabled by default
	}

	return w
}

func Output() (func(string, ...interface{}), io.WriteCloser) {
	exe, _ := os.Executable()
	w := &lumberjack.Logger{
		Filename:   filepath.Join(filepath.Dir(exe), "daemon.log"),
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     7,    //days
		Compress:   true, // disabled by default
	}
	return func(format string, args ...interface{}) {
		header := time.Now().Format("2006-01-02 15:04:05") + " ERROR " + format
		header = header + "\n"
		if w == nil {
			fmt.Printf(header, args...)
			return
		}
		w.Write(S2B(fmt.Sprintf(header, args...)))
	}, w

}

func Upgrade() (func(string, ...interface{}), io.WriteCloser) {
	exe, _ := os.Executable()
	w := &lumberjack.Logger{
		Filename:   filepath.Join(filepath.Dir(exe), "upgrade.log"),
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	return func(format string, args ...interface{}) {
		header := time.Now().Format("2006-01-02 15:04:05") + " ERROR " + format
		header = header + "\n"
		if w == nil {
			fmt.Printf(header, args...)
			return
		}
		w.Write(S2B(fmt.Sprintf(header, args...)))
	}, w
}

func Stdout() (func(string, ...interface{}), io.WriteCloser) {
	exe, _ := os.Executable()
	w := &lumberjack.Logger{
		Filename:   filepath.Join(filepath.Dir(exe), "console.log"),
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	return func(format string, args ...interface{}) {
		header := time.Now().Format("2006-01-02 15:04:05") + " ERROR " + format
		header = header + "\n"
		if w == nil {
			fmt.Printf(header, args...)
			return
		}
		w.Write(S2B(fmt.Sprintf(header, args...)))
	}, w
}

package logger

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"go.uber.org/zap/zapcore"
)

const (
	FormatJson = "json"
	FormatText = "text"

	LevelDebug  = "DEBUG"
	LevelInfo   = "INFO"
	LevelWarn   = "WARN"
	LevelError  = "PTErr"
	LevelDpanic = "DPANIC"
	LevelPanic  = "PTPanic"
	LevelFatal  = "FATAL"
)

type config struct {
	Level      zapcore.Level `ini:"level" yaml:"level" json:"level"`                // 日志输出级别
	Filename   string        `ini:"filename" yaml:"filename" json:"filename"`       // 文件输出位置, 留空则代表不输出到文件
	MaxSize    int           `ini:"maxSize" yaml:"maxSize" json:"maxSize"`          // 单个文件大小, 单位: MiB
	MaxBackups int           `ini:"maxBackups" yaml:"maxBackups" json:"maxBackups"` // 最大文件备份个数
	MaxAge     int           `ini:"maxAge" yaml:"maxAge" json:"maxAge"`             // 日志文件最长留存天数
	Compress   bool          `ini:"compress" yaml:"compress" json:"compress"`       // 备份日志文件是否压缩
	Console    bool          `ini:"console" yaml:"console" json:"console"`          // 是否输出到控制台
	Caller     bool          `ini:"caller" yaml:"caller" json:"caller"`             // 是否打印调用者
	Format     string        `ini:"format" yaml:"format" json:"format"`             // 日志格式化方式
	Color      bool          `ini:"color" yaml:"color" json:"color"`                // 是否显示颜色
	Skip       int           `ini:"skip" yaml:"skip" json:"skip"`                   // 打印代码层级
}

func defaultConfig() *config {
	return &config{
		Level:      zapcore.DebugLevel,
		Filename:   "",
		MaxSize:    100,
		MaxBackups: 100,
		MaxAge:     180,
		Compress:   false,
		//Console:    xEnv.IsDebug(),
		Console: xEnv.IsDebug(),
		Caller:  true,
		Skip:    1,
		Format:  FormatText,
	}
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := defaultConfig()

	tab.Range(func(key string, value lua.LValue) {
		cfg.NewIndex(L, key, value)
	})

	if err := cfg.verify(); err != nil {
		L.RaiseError("logger verify err: %v", err)
		return nil
	}

	return cfg
}

func (c *config) NewIndex(L *lua.LState, key string, val lua.LValue) {

	switch key {

	case "Level":
		c.Level = checkLevel(L, val.String())

	case "Filename":
		c.Filename = val.String()

	case "max_size":
		c.MaxSize = lua.CheckInt(L, val)

	case "max_backups":
		c.MaxBackups = lua.CheckInt(L, val)

	case "max_age":
		c.MaxAge = lua.CheckInt(L, val)

	case "Compress":
		c.Compress = lua.CheckBool(L, val)

	case "console":
		c.Console = lua.CheckBool(L, val)

	case "caller":
		c.Caller = lua.CheckBool(L, val)

	case "skip":
		c.Skip = lua.CheckInt(L, val)

	case "format":
		c.Format = val.String()
	}
}

func (c *config) verify() error {
	if c.Filename == "" && !c.Console {
		return errRequiredOutput
	}
	return nil
}

func checkLevel(L *lua.LState, val string) zapcore.Level {
	var level zapcore.Level

	err := level.UnmarshalText(lua.S2B(val))
	if err != nil {
		L.RaiseError("logger.Level got error: %v", err)
		return 0
	}
	return level
}

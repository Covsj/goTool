package log

import (
	"os"
	"runtime"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	defaultLogger = log.NewJSONLogger(os.Stderr)
)

func Init() {
	InitOpt(level.AllowAll())
}

func InitOpt(opts ...level.Option) {
	defaultLogger = log.With(defaultLogger, "ts", log.DefaultTimestampUTC, "module", "openapi")
	if runtime.GOOS == "darwin" {
		defaultLogger = log.With(defaultLogger, "caller", log.Caller(6))
	} else {
		defaultLogger = log.With(defaultLogger, "caller", log.Caller(8))
	}
	defaultLogger = level.NewFilter(defaultLogger, opts...)
}

func GetDefaultLogger() log.Logger {
	return defaultLogger
}

func SetFilter(opts ...level.Option) {
	defaultLogger = level.NewFilter(defaultLogger, opts...)
}

func Info(kvs ...interface{}) {
	level.Info(defaultLogger).Log(kvs...)
}

func Debug(kvs ...interface{}) {
	level.Debug(defaultLogger).Log(kvs...)
}

func Warn(kvs ...interface{}) {
	level.Warn(defaultLogger).Log(kvs...)
}

func Error(kvs ...interface{}) {
	level.Error(defaultLogger).Log(kvs...)
}

func SetLevel(lv string) {
	switch strings.Trim(strings.ToLower(lv), " ") {
	case "none":
		SetFilter(level.AllowNone())
	case "error":
		SetFilter(level.AllowError())
	case "warn":
		SetFilter(level.AllowWarn())
	case "info":
		SetFilter(level.AllowInfo())
	case "debug":
		SetFilter(level.AllowDebug())
	case "all":
		SetFilter(level.AllowAll())
	default:

	}
}

package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var errLog *logrus.Logger
var once sync.Once

func init() {
	once.Do(func() {
		log = logrus.New()
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "01-02 15:04:05",
		})
		log.SetLevel(logrus.InfoLevel)
		log.SetOutput(os.Stdout)

		errLog = logrus.New()
		errLog.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "01-02 15:04:05",
		})
		errLog.SetLevel(logrus.InfoLevel)
		errLog.SetOutput(os.Stderr)
	})
}

func SetTextFormatter() {
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "01-02 15:04:05",
		ForceColors:     true,
		DisableSorting:  false,
	})

	errLog.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "01-02 15:04:05",
		ForceColors:     true,
		DisableSorting:  false,
	})
}

func addCallerFields(fields map[string]interface{}) logrus.Fields {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	pc, file, line, ok := runtime.Caller(3) // 2 steps up the call stack
	if ok {
		fileName := path.Base(file) // Only the file name
		function := runtime.FuncForPC(pc)
		funcName := ""
		if function != nil {
			parts := strings.Split(function.Name(), ".")
			if len(parts) > 1 {
				funcName = parts[len(parts)-1] // Only the function name
			} else {
				funcName = function.Name()
			}
		}
		// Combine file, function, and line number into a single "caller" attribute
		fields["caller"] = fmt.Sprintf("%s/%s:%d", fileName, funcName, line)
	}
	return fields
}

func argsToFields(args ...interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		var key string
		switch k := args[i].(type) {
		case string:
			key = k
		default:
			key = fmt.Sprintf("%v", k)
		}
		if i+1 >= len(args) {
			fields[key] = ""
		} else {
			fields[key] = args[i+1]
		}
	}
	return fields
}
func logWithFields(logger *logrus.Logger, level logrus.Level, key interface{}, args ...interface{}) {
	fields := addCallerFields(argsToFields(args...))
	switch level {
	case logrus.InfoLevel:
		logger.WithFields(fields).Info(key)
	case logrus.ErrorLevel:
		logger.WithFields(fields).Error(key)
	case logrus.WarnLevel:
		logger.WithFields(fields).Warn(key)
	case logrus.DebugLevel:
		logger.WithFields(fields).Debug(key)
	default:
		logger.WithFields(fields).Info(key)
	}
}

func logWithFieldsFormat(logger *logrus.Logger, level logrus.Level, format string, args ...interface{}) {
	fields := addCallerFields(nil)
	switch level {
	case logrus.InfoLevel:
		logger.WithFields(fields).Infof(format, args...)
	case logrus.ErrorLevel:
		logger.WithFields(fields).Errorf(format, args...)
	case logrus.WarnLevel:
		logger.WithFields(fields).Warnf(format, args...)
	case logrus.DebugLevel:
		logger.WithFields(fields).Debugf(format, args...)
	default:
		logger.WithFields(fields).Infof(format, args...)
	}
}
func Info(key interface{}, args ...interface{}) {
	logWithFields(log, logrus.InfoLevel, key, args...)
}
func InfoF(format string, args ...interface{}) {
	logWithFieldsFormat(log, logrus.InfoLevel, format, args...)
}
func Error(key string, args ...interface{}) {
	logWithFields(errLog, logrus.ErrorLevel, key, args...)
}
func ErrorF(format string, args ...interface{}) {
	logWithFieldsFormat(errLog, logrus.ErrorLevel, format, args...)
}
func Warn(key string, args ...interface{}) {
	logWithFields(log, logrus.WarnLevel, key, args...)
}
func WarnF(format string, args ...interface{}) {
	logWithFieldsFormat(log, logrus.WarnLevel, format, args...)
}
func Debug(key string, args ...interface{}) {
	logWithFields(log, logrus.DebugLevel, key, args...)
}
func DebugF(format string, args ...interface{}) {
	logWithFieldsFormat(log, logrus.DebugLevel, format, args...)
}

type Logger struct {
	*logrus.Logger
	file *os.File
}

func NewLogger(level logrus.Level, format logrus.Formatter, outputFilePath string) (*Logger, error) {
	l := &Logger{
		Logger: logrus.New(),
	}
	l.SetLevel(level)

	if format == nil {
		l.SetFormatter(&logrus.JSONFormatter{})
	} else {
		l.SetFormatter(format)
	}

	if outputFilePath != "" {
		file, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		l.SetOutput(file)
		l.file = file
	} else {
		l.SetOutput(os.Stdout)
	}

	return l, nil
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

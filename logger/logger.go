package logger

import (
	"fmt"
	"time"
)

const (
	LogLevelOff   uint8 = 0
	LogLevelFatal uint8 = 5
	LogLevelError uint8 = 10
	LogLevelWarn  uint8 = 15
	LogLevelInfo  uint8 = 20
	LogLevelDebug uint8 = 25
	LogLevelTrace uint8 = 30
	LogLevelAll   uint8 = 255
)

var LogLevel uint8
var LogDispatcher = func(level uint8, a ...any) {
	fmt.Print(time.Now().Format(time.RFC3339), " - ", LogLevelName(level), " ")
	fmt.Println(a...)
}

func IsEnabled(level uint8) bool {
	return LogLevel >= level
}

func Log(level uint8, a ...any) {
	if IsEnabled(level) {
		LogDispatcher(level, a...)
	}
}

func IsFatalEnabled() bool {
	return IsEnabled(LogLevelFatal)
}

func Fatal(a ...any) {
	Log(LogLevelFatal, a...)
}

func IsErrorEnabled() bool {
	return IsEnabled(LogLevelError)
}

func Error(a ...any) {
	Log(LogLevelError, a...)
}

func IsWarnEnabled() bool {
	return IsEnabled(LogLevelWarn)
}

func Warn(a ...any) {
	Log(LogLevelWarn, a...)
}

func IsInfoEnabled() bool {
	return IsEnabled(LogLevelInfo)
}

func Info(a ...any) {
	Log(LogLevelInfo, a...)
}

func IsDebugEnabled() bool {
	return IsEnabled(LogLevelDebug)
}

func Debug(a ...any) {
	Log(LogLevelDebug, a...)
}

func IsTraceEnabled() bool {
	return IsEnabled(LogLevelTrace)
}

func Trace(a ...any) {
	Log(LogLevelTrace, a...)
}

func IsAllEnabled() bool {
	return IsEnabled(LogLevelAll)
}

func LogLevelName(logLevel uint8) string {
	switch logLevel {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return " INFO"
	case LogLevelWarn:
		return " WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "  LOG"
	}
}

func LogLevelValue(level string) uint8 {
	switch level {
	case "ALL":
		return LogLevelAll
	case "TRACE":
		return LogLevelTrace
	case "DEBUG":
		return LogLevelDebug
	case "INFO":
		return LogLevelInfo
	case "WARN":
		return LogLevelWarn
	case "ERROR":
		return LogLevelError
	case "FATAL":
		return LogLevelFatal
	default:
		return 0
	}
}

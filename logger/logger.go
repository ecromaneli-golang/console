package logger

import (
	"fmt"
	"io"
	"os"
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

var LogLevelStrById = map[uint8]string{
	LogLevelAll:   "ALL",
	LogLevelTrace: "TRACE",
	LogLevelDebug: "DEBUG",
	LogLevelInfo:  "INFO",
	LogLevelWarn:  "WARN",
	LogLevelError: "ERROR",
	LogLevelFatal: "FATAL",
	LogLevelOff:   "OFF",
}

var LogLevelIdByStr = map[string]uint8{
	"ALL":   LogLevelAll,
	"TRACE": LogLevelTrace,
	"DEBUG": LogLevelDebug,
	"INFO":  LogLevelInfo,
	"WARN":  LogLevelWarn,
	"ERROR": LogLevelError,
	"FATAL": LogLevelFatal,
	"OFF":   LogLevelOff,
}

type LogDispatcher func(w io.Writer, name string, level uint8, a ...any)

type logger struct {
	name       string
	dispatcher LogDispatcher
}

var dateFormat = "2006-01-02 15:04:05.000 Z07:00"
var logLevel uint8
var isDefined bool

var globalLogger = New("")

func SetLogLevel(level uint8) {
	logLevel = level
}

func SetLogLevelStr(level string) {
	logLevel = LogLevelIdByStr[level]
}

func SetDateFormat(format string) {
	dateFormat = format
}

func New(name string) *logger {
	return &logger{name: name, dispatcher: DefaultLogDispatcher}
}

func NewCustom(name string, dispatcher func(w io.Writer, name string, level uint8, a ...any)) *logger {
	return &logger{name: name, dispatcher: dispatcher}
}

func GetInstance() *logger {
	return globalLogger
}

func (this *logger) IsEnabled(level uint8) bool {
	return logLevel >= level
}

func (this *logger) Log(level uint8, a ...any) {
	if this.IsEnabled(level) {
		this.dispatcher(os.Stdout, this.name, level, a...)
	}
}

func (this *logger) IsFatalEnabled() bool {
	return this.IsEnabled(LogLevelFatal)
}

func (this *logger) Fatal(a ...any) {
	this.Log(LogLevelFatal, a...)
}

func (this *logger) IsErrorEnabled() bool {
	return this.IsEnabled(LogLevelError)
}

func (this *logger) Error(a ...any) {
	this.Log(LogLevelError, a...)
}

func (this *logger) IsWarnEnabled() bool {
	return this.IsEnabled(LogLevelWarn)
}

func (this *logger) Warn(a ...any) {
	this.Log(LogLevelWarn, a...)
}

func (this *logger) IsInfoEnabled() bool {
	return this.IsEnabled(LogLevelInfo)
}

func (this *logger) Info(a ...any) {
	this.Log(LogLevelInfo, a...)
}

func (this *logger) IsDebugEnabled() bool {
	return this.IsEnabled(LogLevelDebug)
}

func (this *logger) Debug(a ...any) {
	this.Log(LogLevelDebug, a...)
}

func (this *logger) IsTraceEnabled() bool {
	return this.IsEnabled(LogLevelTrace)
}

func (this *logger) Trace(a ...any) {
	this.Log(LogLevelTrace, a...)
}

func (this *logger) IsAllEnabled() bool {
	return this.IsEnabled(LogLevelAll)
}

func DefaultLogDispatcher(w io.Writer, name string, level uint8, a ...any) {
	if name != "" {
		name = " " + name + ":"
	}

	levelStr := LogLevelStrById[level]
	if len(levelStr) == 4 {
		levelStr += " "
	}

	fmt.Fprint(w, time.Now().Format(dateFormat)+" - "+levelStr+name+" ")
	fmt.Fprintln(w, a...)
}

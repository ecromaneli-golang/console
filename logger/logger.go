package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level uint8

const (
	LevelOff   Level = 0
	LevelFatal Level = 5
	LevelError Level = 10
	LevelWarn  Level = 15
	LevelInfo  Level = 20
	LevelDebug Level = 25
	LevelTrace Level = 30
	LevelAll   Level = 255
)

func (l *Level) String() string {
	return stringByLevel[*l]
}

func LevelFromString(level string) Level {
	return levelByStr[level]
}

var stringByLevel = map[Level]string{
	LevelAll:   "ALL",
	LevelTrace: "TRACE",
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
	LevelOff:   "OFF",
}

var levelByStr = map[string]Level{
	"ALL":   LevelAll,
	"TRACE": LevelTrace,
	"DEBUG": LevelDebug,
	"INFO":  LevelInfo,
	"WARN":  LevelWarn,
	"ERROR": LevelError,
	"FATAL": LevelFatal,
	"OFF":   LevelOff,
}

type LogDispatcher func(w io.Writer, dateFormat string, name string, level Level, a ...any)

type Logger struct {
	name       string
	dispatcher LogDispatcher
	writer     io.Writer
	logLevel   Level
	dateFormat string
}

var (
	DefaultDateFormat           = "2006-01-02 15:04:05.000 Z07:00"
	DefaultWriter     io.Writer = os.Stdout
	DefaultDispatcher           = DefaultLogDispatcher
	DefaultLogLevel             = LevelInfo
	globalLogger      *Logger
)

func GetInstance() *Logger {
	if globalLogger == nil {
		globalLogger = New("")
	}
	return globalLogger
}

func SetDefaultDateFormat(format string) error {
	DefaultDateFormat = format
	return nil
}

func SetDefaultOutput(writer io.Writer) {
	DefaultWriter = writer
}

func SetDefaultLogDispatcher(dispatcher LogDispatcher) {
	DefaultDispatcher = dispatcher
}

func SetDefaultLogLevel(l Level) {
	DefaultLogLevel = l
}

func New(name string) *Logger {
	return &Logger{
		name:       name,
		dispatcher: DefaultDispatcher,
		writer:     DefaultWriter,
		logLevel:   DefaultLogLevel,
		dateFormat: DefaultDateFormat,
	}
}

func (l *Logger) SetLogLevel(lv Level) {
	l.logLevel = lv
}

func (l *Logger) SetLogLevelStr(levelStr string) {
	l.logLevel = LevelFromString(levelStr)
}

func (l *Logger) SetDateFormat(format string) error {
	l.dateFormat = format
	return nil
}

func (l *Logger) SetLogDispatcher(dispatcher LogDispatcher) {
	l.dispatcher = dispatcher
}

func (l *Logger) SetOutput(writer io.Writer) {
	l.writer = writer
}

func (l *Logger) IsEnabled(lv Level) bool {
	return l.logLevel >= lv
}

func (l *Logger) Log(lv Level, a ...any) {
	if l.IsEnabled(lv) {
		l.dispatcher(l.writer, l.dateFormat, l.name, lv, a...)
	}
}

func (l *Logger) Fatal(a ...any) {
	l.Log(LevelFatal, a...)
}

func (l *Logger) Error(a ...any) {
	l.Log(LevelError, a...)
}

func (l *Logger) Warn(a ...any) {
	l.Log(LevelWarn, a...)
}

func (l *Logger) Info(a ...any) {
	l.Log(LevelInfo, a...)
}

func (l *Logger) Debug(a ...any) {
	l.Log(LevelDebug, a...)
}

func (l *Logger) Trace(a ...any) {
	l.Log(LevelTrace, a...)
}

func DefaultLogDispatcher(w io.Writer, dateFormat string, name string, l Level, a ...any) {
	message := ""

	if dateFormat != "" {
		message = time.Now().Format(dateFormat) + " - "
	}

	levelStr := l.String()
	if len(levelStr) == 4 {
		levelStr += " "
	}
	message += levelStr

	if name != "" {
		message += " " + name + ":"
	}

	fmt.Fprintln(w, message, fmt.Sprint(a...))
}

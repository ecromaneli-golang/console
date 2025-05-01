package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ecromaneli-golang/console/logger/async"
)

// Level represents the severity of a log message.
type Level uint8

const (
	// LevelOff disables all logging.
	LevelOff Level = 0
	// LevelFatal is for critical errors that cause application failure.
	LevelFatal Level = 5
	// LevelError is for runtime errors that don't cause application failure.
	LevelError Level = 10
	// LevelWarn is for potentially harmful situations.
	LevelWarn Level = 15
	// LevelInfo is for general information messages.
	LevelInfo Level = 20
	// LevelDebug is for detailed debugging information.
	LevelDebug Level = 25
	// LevelTrace is for the most fine-grained information.
	LevelTrace Level = 30
	// LevelAll enables all possible logging levels.
	LevelAll Level = 255
)

// String returns the string representation of the log level.
func (l *Level) String() string {
	return stringByLevel[*l]
}

// LevelFromString converts a string representation to a Level.
//
// It returns the corresponding Level for the given string value.
func LevelFromString(level string) Level {
	return levelByStr[strings.ToUpper(level)]
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

// LogDispatcher is a function type that handles formatting and writing log messages.
type LogDispatcher func(w io.Writer, dateFormat string, name string, level Level, a ...any)

// Logger provides methods for logging messages at different levels.
type Logger struct {
	name       string
	dispatcher LogDispatcher
	writer     io.Writer
	logLevel   Level
	dateFormat string
}

var (
	// DefaultDateFormat is the default format for timestamps in log messages.
	DefaultDateFormat = "2006-01-02 15:04:05.000 Z07:00"
	// DefaultWriter is the default output destination for log messages.
	DefaultWriter io.Writer = os.Stdout
	// DefaultDispatcher is the default function used to format and write log messages.
	DefaultDispatcher = DefaultLogDispatcher
	// DefaultLogLevel is the default level at which messages are logged.
	DefaultLogLevel = LevelInfo
	globalLogger    *Logger
)

// GetInstance returns the global logger instance, creating it if it doesn't exist.
//
// This provides a singleton pattern for accessing a shared logger.
func GetInstance() *Logger {
	if globalLogger == nil {
		globalLogger = New("")
	}
	return globalLogger
}

// SetDefaultDateFormat sets the default date format for new logger instances.
//
// The format should be compatible with Go's time.Format function.
func SetDefaultDateFormat(format string) error {
	DefaultDateFormat = format
	return nil
}

// SetDefaultOutput sets the default writer for new logger instances.
//
// The writer is where log messages will be written.
func SetDefaultOutput(writer io.Writer) {
	DefaultWriter = writer
}

// SetDefaultLogDispatcher sets the default dispatcher function for new logger instances.
//
// The dispatcher controls how log messages are formatted and written.
func SetDefaultLogDispatcher(dispatcher LogDispatcher) {
	DefaultDispatcher = dispatcher
}

// SetDefaultLogLevel sets the default minimum level for new logger instances.
//
// Messages below this level will not be logged.
func SetDefaultLogLevel(l Level) {
	DefaultLogLevel = l
}

// SetDefaultLogLevelStr sets the default minimum level using a string representation.
//
// It converts the string to the corresponding Level and sets it as the default.
func SetDefaultLogLevelStr(levelStr string) {
	DefaultLogLevel = LevelFromString(levelStr)
}

// New creates a new logger with the given name and default settings.
//
// The name is included in log messages to identify their source.
func New(name string) *Logger {
	return &Logger{
		name:       name,
		dispatcher: DefaultDispatcher,
		writer:     DefaultWriter,
		logLevel:   DefaultLogLevel,
		dateFormat: DefaultDateFormat,
	}
}

// SetLogLevel sets the minimum level at which messages will be logged.
//
// Messages below this level will not be logged.
func (l *Logger) SetLogLevel(lv Level) {
	l.logLevel = lv
}

// SetLogLevelStr sets the minimum log level using a string representation.
//
// It converts the string to the corresponding Level and sets it.
func (l *Logger) SetLogLevelStr(levelStr string) {
	l.logLevel = LevelFromString(levelStr)
}

// SetDateFormat sets the date format used in log messages.
//
// The format should be compatible with Go's time.Format function.
func (l *Logger) SetDateFormat(format string) error {
	l.dateFormat = format
	return nil
}

// SetLogDispatcher sets the dispatcher function for this logger.
//
// The dispatcher controls how log messages are formatted and written.
func (l *Logger) SetLogDispatcher(dispatcher LogDispatcher) {
	l.dispatcher = dispatcher
}

// SetOutput sets the output where log messages will be written.
func (l *Logger) SetOutput(writer io.Writer) {
	l.writer = writer
}

// SetAsyncOutput sets the output to an asynchronous writer.
// Shorthand for l.SetOutput(async.NewAsyncWriter([...])).
//
// The bufferSize determines the number of pending writes that can be queued before blocking.
// If the buffer is full, the log message will be written directly to the target writer.
func (l *Logger) SetAsyncOutput(writer io.Writer, bufferSize int) {
	if _, ok := l.writer.(*async.AsyncWriter); ok {
		l.writer = writer
	} else {
		l.writer = async.NewAsyncWriter(writer, bufferSize)
	}
}

// SetAsync sets the current output to an asynchronous writer.
// Shorthand for l.SetOutput(async.NewAsyncWriter(currentOutput, bufferSize)).
//
// The bufferSize determines the number of pending writes that can be queued before blocking.
// If the buffer is full, the log message will be written directly to the target writer.
func (l *Logger) SetAsync(bufferSize int) {
	if _, ok := l.writer.(*async.AsyncWriter); !ok {
		l.writer = async.NewAsyncWriter(l.writer, bufferSize)
	}
}

// SetSync sets the current output to synchronous mode.
func (l *Logger) SetSync() {
	if asyncWriter, ok := l.writer.(*async.AsyncWriter); ok {
		asyncWriter.Flush()
		l.writer = asyncWriter.Target()
	}
}

// IsEnabled returns true if the given level is enabled for logging.
//
// A level is enabled if it is greater than or equal to the logger's level.
func (l *Logger) IsEnabled(lv Level) bool {
	return l.logLevel >= lv
}

// IsFatalEnabled returns true if fatal level messages will be logged.
//
// This is a convenience method equivalent to IsEnabled(LevelFatal).
func (l *Logger) IsFatalEnabled() bool {
	return l.IsEnabled(LevelFatal)
}

// IsErrorEnabled returns true if error level messages will be logged.
//
// This is a convenience method equivalent to IsEnabled(LevelError).
func (l *Logger) IsErrorEnabled() bool {
	return l.IsEnabled(LevelError)
}

// IsWarnEnabled returns true if warning level messages will be logged.
//
// This is a convenience method equivalent to IsEnabled(LevelWarn).
func (l *Logger) IsWarnEnabled() bool {
	return l.IsEnabled(LevelWarn)
}

// IsInfoEnabled returns true if info level messages will be logged.
//
// This is a convenience method equivalent to IsEnabled(LevelInfo).
func (l *Logger) IsInfoEnabled() bool {
	return l.IsEnabled(LevelInfo)
}

// IsDebugEnabled returns true if debug level messages will be logged.
//
// This is a convenience method equivalent to IsEnabled(LevelDebug).
func (l *Logger) IsDebugEnabled() bool {
	return l.IsEnabled(LevelDebug)
}

// IsTraceEnabled returns true if trace level messages will be logged.
//
// This is a convenience method equivalent to IsEnabled(LevelTrace).
func (l *Logger) IsTraceEnabled() bool {
	return l.IsEnabled(LevelTrace)
}

// Log logs a message at the specified level.
//
// If the level is enabled, the message is passed to the dispatcher.
func (l *Logger) Log(lv Level, a ...any) {
	if l.IsEnabled(lv) {
		l.dispatcher(l.writer, l.dateFormat, l.name, lv, a...)
	}
}

// Fatal logs a message at the fatal level.
//
// This should be used for critical errors that cause application failure.
func (l *Logger) Fatal(a ...any) {
	l.Log(LevelFatal, a...)
}

// Error logs a message at the error level.
//
// This should be used for runtime errors that don't cause application failure.
func (l *Logger) Error(a ...any) {
	l.Log(LevelError, a...)
}

// Warn logs a message at the warning level.
//
// This should be used for potentially harmful situations.
func (l *Logger) Warn(a ...any) {
	l.Log(LevelWarn, a...)
}

// Info logs a message at the info level.
//
// This should be used for general information messages.
func (l *Logger) Info(a ...any) {
	l.Log(LevelInfo, a...)
}

// Debug logs a message at the debug level.
//
// This should be used for detailed debugging information.
func (l *Logger) Debug(a ...any) {
	l.Log(LevelDebug, a...)
}

// Trace logs a message at the trace level.
//
// This should be used for the most fine-grained information.
func (l *Logger) Trace(a ...any) {
	l.Log(LevelTrace, a...)
}

// Flush waits for all pending writes to complete.
// It is only necessary if the logger is using an asynchronous writer.
// Shorthand for logger.Output().Flush().
func (l *Logger) Flush() {
	if asyncWriter, ok := l.writer.(*async.AsyncWriter); ok {
		asyncWriter.Flush()
	}
}

// Name returns the name of the logger.
//
// The name is included in log messages to identify their source.
func (l *Logger) Name() string {
	return l.name
}

// Dispatcher returns the current log dispatcher function.
//
// The dispatcher controls how log messages are formatted and written.
func (l *Logger) Dispatcher() LogDispatcher {
	return l.dispatcher
}

// Output returns the current writer where log messages are written.
//
// This is the destination for all log messages.
func (l *Logger) Output() io.Writer {
	return l.writer
}

// LogLevel returns the minimum level at which messages will be logged.
//
// Messages below this level will not be logged.
func (l *Logger) LogLevel() Level {
	return l.logLevel
}

// DateFormat returns the date format used in log messages.
//
// The format is compatible with Go's time.Format function.
func (l *Logger) DateFormat() string {
	return l.dateFormat
}

// DefaultLogDispatcher is the default function for formatting and writing log messages.
//
// It formats the message with a timestamp, log level, name, and the message content.
func DefaultLogDispatcher(w io.Writer, dateFormat string, name string, l Level, a ...any) {
	estSize := 6

	if dateFormat != "" {
		dateFormat = time.Now().Format(dateFormat)
		estSize += len(dateFormat) + 3 // 3 for " - "
	}

	if name != "" {
		estSize += len(name) + 2 // 2 for " :"
	}

	message := fmt.Sprintln(a...)

	estSize += len(message)

	var builder strings.Builder
	builder.Grow(estSize)

	// Add the timestamp if a date format is provided
	if dateFormat != "" {
		builder.WriteString(dateFormat)
		builder.WriteString(" - ")
	}

	// Add the log level
	levelStr := l.String()
	builder.WriteString(levelStr)
	if len(levelStr) == 4 {
		builder.WriteByte(' ')
	}

	// Add the logger name if provided
	if name != "" {
		builder.WriteByte(' ')
		builder.WriteString(name)
		builder.WriteByte(':')
	}

	// Add a space before the message
	builder.WriteByte(' ')

	// Add the log message
	builder.WriteString(message)

	// Write the final message to the writer
	fmt.Fprint(w, builder.String())
}

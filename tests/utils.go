package tests

import (
	"fmt"
	"io"
	"testing"

	"github.com/ecromaneli-golang/console/logger"
)

type LogCounter map[logger.Level][]Logger

type Logger struct {
	Name       string
	Writer     io.Writer
	LogLevel   logger.Level
	DateFormat string
	Message    string
}

func (c *LogCounter) GetTotal() int {
	var total int

	for _, logs := range *c {
		total += len(logs)
	}

	return total
}

func NewCounterDispatcher() (logger.LogDispatcher, LogCounter) {
	counter := make(LogCounter)

	return func(w io.Writer, dateFormat string, name string, level logger.Level, a ...any) {
		// Log the message using the default dispatcher
		logger.DefaultLogDispatcher(w, dateFormat, name, level, a...)

		// Ensure the slice for the log level is initialized
		if _, exists := counter[level]; !exists {
			counter[level] = []Logger{}
		}

		// Store the log data in the counter
		logEntry := Logger{
			Name:       name,
			Writer:     w,
			LogLevel:   level,
			DateFormat: dateFormat,
			Message:    fmt.Sprint(a...),
		}

		counter[level] = append(counter[level], logEntry)
	}, counter
}

func AssertEquals(t *testing.T, expected any, current any) {
	if expected != current {
		t.Errorf("\n\nExpected: %v\nCurrent: %v\n\n", expected, current)
	}
}

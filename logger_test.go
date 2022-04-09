package logger

import (
	"testing"

	"github.com/ecromaneli-golang/console/logger"
)

type log struct {
	level uint8
	a     []any
}

var stdout chan log = make(chan log, 16)

func FakeDispatcher(level uint8, a ...any) {
	stdout <- log{level: level, a: a}
}

func TestShouldLogFatalAndError(t *testing.T) {
	// Given
	logger.LogLevel = logger.LogLevelError
	logger.LogDispatcher = FakeDispatcher

	// When
	logger.Fatal("Fatal should log")
	logger.Error("Error should log")
	logger.Debug("Debug should not log")
	close(stdout)

	// Then
	logFatal := false
	logError := false

	for log := range stdout {
		if log.level == logger.LogLevelDebug {
			t.FailNow()
		}

		if log.level == logger.LogLevelFatal {
			logFatal = true
		}

		if log.level == logger.LogLevelError {
			logError = true
		}
	}

	if !logFatal || !logError {
		t.FailNow()
	}
}

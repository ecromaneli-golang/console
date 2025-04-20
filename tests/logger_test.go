package tests

import (
	"testing"

	"github.com/ecromaneli-golang/console/logger"
)

func TestShouldLogFatalAndError(t *testing.T) {
	// Given
	logger.SetDefaultLogLevel(logger.LevelError)

	dispatcher, counter := NewCounterDispatcher()
	log := logger.New("test")
	log.SetLogDispatcher(dispatcher)

	// When
	log.Fatal("1")
	log.Error("2")
	log.Debug("ignored")

	// Then
	AssertEquals(t, 1, len(counter[logger.LevelFatal]))
	AssertEquals(t, 1, len(counter[logger.LevelError]))
	AssertEquals(t, 2, counter.GetTotal())
}

func TestShouldUseGlobalInstance(t *testing.T) {
	// Given
	logger.SetDefaultLogLevel(logger.LevelAll)
	log := logger.GetInstance()

	// When
	log.Fatal("Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	log.Error("Phasellus eu odio libero. Curabitur sed elit dictum")
	log.Warn("Sed ligula mauris, rutrum ac ipsum eget")
	log.Info("Duis non finibus erat. In consectetur facilisis purus ac rhoncus")
	log.Debug("Class aptent taciti sociosqu ad litora torquent")
	log.Trace("Sed tincidunt egestas dolor, nec tincidunt tortor accumsan ac")

	// Then no panic
}

func TestShouldNotPrintDate(t *testing.T) {
	// Given
	logger.SetDefaultLogLevel(logger.LevelAll)

	dispatcher, counter := NewCounterDispatcher()
	log := logger.New("test")
	log.SetLogDispatcher(dispatcher)
	log.SetDateFormat("")

	// When
	log.Fatal("Lorem ipsum dolor sit amet, consectetur adipiscing elit")

	// Then
	AssertEquals(t, "", counter[logger.LevelFatal][0].DateFormat)
}

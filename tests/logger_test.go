package tests

import (
	"sync"
	"testing"

	"github.com/ecromaneli-golang/console/logger"
)

var mux sync.Mutex

func TestShouldLogFatalAndError(t *testing.T) {
	mux.Lock()
	defer mux.Unlock()

	// Given
	logger.SetLogLevel(logger.LogLevelError)

	dispatcher, counter := NewCounterDispatcher()
	log := logger.NewCustom("test", dispatcher)

	// When
	log.Fatal("1")
	log.Error("2")
	log.Debug("ignored")

	// Then
	assertEquals(t, 1, counter[logger.LogLevelFatal])
	assertEquals(t, 1, counter[logger.LogLevelError])
	assertEquals(t, 2, counter.GetTotal())
}

func TestShouldPrintAllLogs(t *testing.T) {
	mux.Lock()
	defer mux.Unlock()

	// Given
	logger.SetLogLevel(logger.LogLevelAll)

	dispatcher, counter := NewCounterDispatcher()
	log := logger.NewCustom("test", dispatcher)

	// When
	log.Fatal("Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	log.Error("Phasellus eu odio libero. Curabitur sed elit dictum")
	log.Warn("Sed ligula mauris, rutrum ac ipsum eget")
	log.Info("Duis non finibus erat. In consectetur facilisis purus ac rhoncus")
	log.Debug("Class aptent taciti sociosqu ad litora torquent")
	log.Trace("Sed tincidunt egestas dolor, nec tincidunt tortor accumsan ac")

	// Then
	for _, v := range counter {
		assertEquals(t, 1, v)
	}

	assertEquals(t, 6, counter.GetTotal())
}

func TestShouldUseGlobalInstance(t *testing.T) {
	mux.Lock()
	defer mux.Unlock()

	// Given
	logger.SetLogLevel(logger.LogLevelAll)
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

func assertEquals(t *testing.T, expected any, current any) {
	if expected != current {
		t.Errorf("\n\nExpected: %v\nCurrent: %v\n\n", expected, current)
	}
}

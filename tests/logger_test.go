package tests

import (
	"bytes"
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

func TestShouldParseStringToLevel(t *testing.T) {
	// Given
	dispatcher, counter := NewCounterDispatcher()
	log := logger.New("test")
	log.SetLogDispatcher(dispatcher)

	// When
	log.SetLogLevelStr("fAtAL")
	log.Fatal("Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	log.Error("Phasellus eu odio libero. Curabitur sed elit dictum")

	// Then
	AssertEquals(t, 1, len(counter[logger.LevelFatal]))
	AssertEquals(t, 0, len(counter[logger.LevelError]))
}

func TestShouldNotTriggerDisabledCode(t *testing.T) {
	// Given
	dispatcher, _ := NewCounterDispatcher()
	log := logger.New("test")
	log.SetLogDispatcher(dispatcher)
	log.SetLogLevelStr("fAtAL")

	// When

	if log.IsTraceEnabled() {
		panic("Should not be enabled")
	}
	if log.IsDebugEnabled() {
		panic("Should not be enabled")
	}
	if log.IsInfoEnabled() {
		panic("Should not be enabled")
	}
	if log.IsWarnEnabled() {
		panic("Should not be enabled")
	}
	if log.IsErrorEnabled() {
		panic("Should not be enabled")
	}
	if !log.IsFatalEnabled() {
		panic("Should be enabled")
	}

	// Then no panic
}

// Fixed by v1.3.1 - After change to fmt.Sprint(a...), the output was not concat with spaces anymore
func TestShouldUseDefaultFormatting(t *testing.T) {
	// Given
	var output bytes.Buffer

	log := logger.New("AnyName")
	log.SetDateFormat("AnyDate")
	log.SetOutput(&output)

	// When
	log.Warn("split", "test", 1, 2, 3)

	// Then
	AssertEquals(t, "AnyDate - WARN  AnyName: split test 1 2 3\n", output.String())
}

func TestShouldLogAsync(t *testing.T) {
	// Given
	const anyMessage = "any message"
	var output bytes.Buffer

	log := logger.New("AnyName")
	log.SetLogLevelStr("fAtAL")
	log.SetLogDispatcher(UnformattedDispatcher)
	log.SetAsyncOutput(&output, 1)

	// When
	log.Fatal(anyMessage)

	// Then
	AssertEquals(t, "", output.String())
	log.Flush()
	AssertEquals(t, anyMessage, output.String())
}

func TestShouldLogAsyncWithBuffer(t *testing.T) {
	// Given
	var output bytes.Buffer
	log := logger.New("AnyName")
	log.SetLogLevelStr("fatal")
	log.SetAsyncOutput(&output, 10)
	log.SetLogDispatcher(UnformattedDispatcher)

	// When
	log.Fatal("1")
	log.Fatal("2")
	log.Fatal("3")
	log.Fatal("4")
	log.Fatal("5")
	log.Fatal("6")
	log.Fatal("7")
	log.Fatal("8")
	log.Fatal("9")
	log.Fatal("0")
	log.Error("Will not be printed")
	log.Flush()

	// Then
	AssertEquals(t, "1234567890", output.String())
}

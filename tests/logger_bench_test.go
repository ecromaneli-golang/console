package tests

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/ecromaneli-golang/console/logger"
	"github.com/ecromaneli-golang/console/logger/async"
)

func BenchmarkDefaultLogDispatcher(b *testing.B) {
	// Given
	var output bytes.Buffer
	dateFormat := "2006-01-02 15:04:05.000 Z07:00"
	name := "BenchmarkLogger"
	level := logger.LevelInfo
	message := "This is a benchmark test message"

	// When
	for i := 0; b.Loop(); i++ {
		logger.DefaultLogDispatcher(&output, dateFormat, name, level, strconv.Itoa(i)+" - "+message)
	}
}

func BenchmarkAsyncLogDispatcher(b *testing.B) {
	// Given
	var output bytes.Buffer
	asyncWriter := async.NewAsyncWriter(&output, 0)
	dateFormat := "2006-01-02 15:04:05.000 Z07:00"
	name := "BenchmarkLogger"
	level := logger.LevelInfo
	message := "This is a benchmark test message"

	// When
	for i := 0; b.Loop(); i++ {
		logger.DefaultLogDispatcher(asyncWriter, dateFormat, name, level, strconv.Itoa(i)+" - "+message)
	}
}

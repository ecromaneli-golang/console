package tests

import (
	"bytes"
	"testing"

	"github.com/ecromaneli-golang/console/logger"
)

func BenchmarkDefaultLogDispatcher(b *testing.B) {
	// Given
	var output bytes.Buffer
	dateFormat := "2006-01-02 15:04:05.000 Z07:00"
	name := "BenchmarkLogger"
	level := logger.LevelInfo
	message := "This is a benchmark test message"

	// When
	for i := 0; i < b.N; i++ {
		logger.DefaultLogDispatcher(&output, dateFormat, name, level, message)
	}
}

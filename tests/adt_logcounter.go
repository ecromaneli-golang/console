package tests

import (
	"io"

	"github.com/ecromaneli-golang/console/logger"
)

type LogCounter map[uint8]int

func (this *LogCounter) GetTotal() int {
	var total int

	for _, count := range *this {
		total += count
	}

	return total
}

func NewCounterDispatcher() (logger.LogDispatcher, LogCounter) {
	counter := make(LogCounter)

	return func(w io.Writer, name string, level uint8, a ...any) {
		logger.DefaultLogDispatcher(w, name, level, a...)
		counter[level]++
	}, counter
}

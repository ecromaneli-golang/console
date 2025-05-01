package async

import (
	"io"
	"sync"
)

// AsyncWriter is an io.Writer that processes writes asynchronously.
type AsyncWriter struct {
	target      io.Writer      // The actual destination writer
	dataCh      chan []byte    // Channel for the data to be written
	done        chan any       // Signal channel for shutdown
	wg          sync.WaitGroup // To ensure all writes complete
	closed      bool
	forceClosed bool
	bufSize     int
}

// NewAsyncWriter creates a new AsyncWriter that writes to the target writer asynchronously.
// bufferSize determines the number of pending writes that can be queued before blocking.
func NewAsyncWriter(target io.Writer, bufferSize int) *AsyncWriter {
	if bufferSize <= 0 {
		bufferSize = 100 // Default buffer size
	}

	w := &AsyncWriter{
		target:  target,
		dataCh:  make(chan []byte, bufferSize),
		done:    make(chan any),
		bufSize: bufferSize,
	}

	w.wg.Add(1)
	go w.processWrites()

	return w
}

// Write implements io.Writer and sends data to be written asynchronously.
func (w *AsyncWriter) Write(p []byte) (n int, err error) {
	if w.closed {
		return 0, io.ErrClosedPipe
	}

	// Make a copy of the data since p may be reused by the caller
	length := len(p)
	data := make([]byte, length)
	copy(data, p)

	// Try to send without blocking
	select {
	case w.dataCh <- data:
		return length, nil
	default:
		// Channel is full, write directly to target
		return w.target.Write(p)
	}
}

// processWrites reads from the channel and writes to the target writer.
func (w *AsyncWriter) processWrites() {
	defer w.wg.Done()

	for {
		select {
		case data := <-w.dataCh:
			w.target.Write(data)
		case <-w.done:
			// Process any remaining writes
			close(w.dataCh)
			if !w.forceClosed {
				for data := range w.dataCh {
					w.target.Write(data)
				}
			}
			return
		}
	}
}

// Flush waits for all pending writes to complete.
func (w *AsyncWriter) Flush() {
	if !w.closed {
		w.closed = true
		close(w.done)
		w.wg.Wait()
	}
}

// Close closes the AsyncWriter without waiting for all pending writes to complete.
func (w *AsyncWriter) Close() {
	if !w.closed {
		w.closed = true
		w.forceClosed = true
		close(w.done)
	}
}

// Target returns the underlying target writer.
func (w *AsyncWriter) Target() io.Writer {
	return w.target
}

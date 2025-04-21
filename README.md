# Logger Library

A lightweight and flexible logging library for Go, designed to support multiple log levels, custom dispatchers, and configurable outputs. This library is ideal for applications that require structured and customizable logging.

[![Go Reference](https://pkg.go.dev/badge/github.com/ecromaneli-golang/console.svg)](https://pkg.go.dev/github.com/ecromaneli-golang/console)
[![Go Report Card](https://goreportcard.com/badge/github.com/ecromaneli-golang/console)](https://goreportcard.com/report/github.com/ecromaneli-golang/console)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Log Levels**: Supports multiple log levels (`Fatal`, `Error`, `Warn`, `Info`, `Debug`, `Trace`, etc.).
- **Custom Dispatchers**: Define how log messages are formatted and written.
- **Global and Instance-Based Loggers**: Use a global logger or create multiple logger instances.
- **Configurable Outputs**: Write logs to `os.Stdout`, files, or any `io.Writer`.
- **Thread-Safe**: Designed to work in concurrent environments.

## Installation

To use this library, add it to your project using `go get`:

```bash
go get github.com/ecromaneli-golang/console/logger
```

## Usage

### Basic Example

```go
package main

import (
	"github.com/ecromaneli-golang/console/logger"
)

func main() {
	log := logger.New("MyApp")
	log.SetLogLevel(logger.LevelInfo)

	log.Info("Application started")
	log.Warn("This is a warning")
	log.Error("An error occurred")
}
```

### Output

```
2025-04-20 15:04:05.000 Z07:00 - INFO  MyApp: Application started
2025-04-20 15:04:05.000 Z07:00 - WARN  MyApp: This is a warning
2025-04-20 15:04:05.000 Z07:00 - ERROR MyApp: An error occurred
```

### Using the Global Logger

```go
package main

import (
	"github.com/ecromaneli-golang/console/logger"
)

func main() {
	log := logger.GetInstance()
	log.SetLogLevelStr("DEBUG")

	log.Debug("Debugging details")
	log.Fatal("Critical failure")
}
```

### Output

```
2025-04-20 15:04:05.000 Z07:00 - DEBUG Debugging details
2025-04-20 15:04:05.000 Z07:00 - FATAL Critical failure
```

## Configuration

### Setting Default Values

You can configure default values for all new loggers:

```go
logger.SetDefaultLogLevel(logger.LevelWarn)
logger.SetDefaultDateFormat("2006-01-02")
logger.SetDefaultOutput(os.Stderr)
logger.SetDefaultLogDispatcher(func(w io.Writer, dateFormat string, name string, level logger.Level, a ...any) {
	fmt.Fprintf(w, "[CUSTOM] %s - %s: %s\n", dateFormat, name, fmt.Sprint(a...))
})
```

### Customizing a Logger Instance

Each logger instance can be customized independently:

```go
log := logger.New("CustomLogger")
log.SetLogLevel(logger.LevelTrace)
log.SetDateFormat("15:04:05")
log.SetOutput(os.Stdout)
log.SetLogDispatcher(func(w io.Writer, dateFormat string, name string, level logger.Level, a ...any) {
	fmt.Fprintf(w, "[%s] %s - %s: %s\n", level.String(), dateFormat, name, fmt.Sprint(a...))
})
```

### Output

```
[TRACE] 15:04:05 - CustomLogger: Custom trace message
```

## Log Levels

The library supports the following log levels:

| LevelStr    | Level       | Value | Description                     |
|-------------|-------------|-------|---------------------------------|
| `OFF`       | `LevelOff`  | `0`   | Disables all logging.           |
| `FATAL`     | `LevelFatal`| `5`   | Logs critical errors.           |
| `ERROR`     | `LevelError`| `10`  | Logs errors.                    |
| `WARN`      | `LevelWarn` | `15`  | Logs warnings.                  |
| `INFO`      | `LevelInfo` | `20`  | Logs informational messages.    |
| `DEBUG`     | `LevelDebug`| `25`  | Logs debug messages.            |
| `TRACE`     | `LevelTrace`| `30`  | Logs trace messages.            |
| `ALL`       | `LevelAll`  | `255` | Enables all logging levels.     |

## Advanced Usage

### Custom Log Dispatcher

You can define a custom dispatcher to control how logs are formatted and written:

```go
log := logger.New("AdvancedLogger")
log.SetLogDispatcher(func(w io.Writer, dateFormat string, name string, level logger.Level, a ...any) {
	fmt.Fprintf(w, "[%s] %s - %s: %s\n", level.String(), dateFormat, name, fmt.Sprint(a...))
})
log.Info("Custom log format")
```

### Output

```
[INFO] 2025-04-20 15:04:05.000 Z07:00 - AdvancedLogger: Custom log format
```

### Logging Without a Date

To disable the date in logs, set an empty date format:

```go
log := logger.New("NoDateLogger")
log.SetDateFormat("")
log.Info("This log has no date")
```

### Output

```
INFO NoDateLogger: This log has no date
```

## Testing

The library includes utilities for testing loggers, such as `NewCounterDispatcher` to count log messages by level.

### Example Test

```go
package tests

import (
	"testing"

	"github.com/ecromaneli-golang/console/logger"
)

func TestLogger(t *testing.T) {
	dispatcher, counter := NewCounterDispatcher()
	log := logger.New("TestLogger")
	log.SetLogDispatcher(dispatcher)

	log.Info("Info message")
	log.Warn("Warning message")

	if counter.GetTotal() != 2 {
		t.Errorf("Expected 2 logs, got %d", counter.GetTotal())
	}
}
```

### Output

```
2025-04-20 15:04:05.000 Z07:00 - INFO  TestLogger: Info message
2025-04-20 15:04:05.000 Z07:00 - WARN  TestLogger: Warning message
```

## Author

- **Author**: Emerson C. Romaneli
- **GitHub**: [ecromaneli](https://github.com/ecromaneli)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Commit your changes with clear messages.
4. Submit a pull request.

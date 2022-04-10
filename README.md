# Logger
Go Logger Basic Implementation

# Basic Usage

```go
    package main

    import "github.com/ecromaneli-golang/console/logger"

    log := logger.New("logname")

    func main() {
        // Configuring global log level
        logger.SetLogLevel(logger.LogLevelDebug)
        /* OR */
        logger.SetLogLevelStr("DEBUG")

        // Usage
        log.Info("As simple as that")

        // Will print
        "2022-04-08 22:58:45.735 -03:00 - INFO logname: As simple as that"
    }

```
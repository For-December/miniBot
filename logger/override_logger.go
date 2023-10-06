package logger

import (
	"log"
)

type OverrideLogger struct {
	logger *log.Logger
}

// New creates a new FileLogger.

// Debug logs a message at DebugLevel. The message includes any fields passed
func (f OverrideLogger) Debug(v ...interface{}) {}

// Info logs a message at InfoLevel. The message includes any fields passed
func (f OverrideLogger) Info(v ...interface{}) {}

// Warn logs a message at WarnLevel. The message includes any fields passed
func (f OverrideLogger) Warn(v ...interface{}) {}

// Error logs a message at ErrorLevel. The message includes any fields passed
func (f OverrideLogger) Error(v ...interface{}) {}

// Debugf logs a message at DebugLevel. The message includes any fields passed
func (f OverrideLogger) Debugf(format string, v ...interface{}) {}

// Infof logs a message at InfoLevel. The message includes any fields passed
func (f OverrideLogger) Infof(format string, v ...interface{}) {}

// Warnf logs a message at WarnLevel. The message includes any fields passed
func (f OverrideLogger) Warnf(format string, v ...interface{}) {}

// Errorf logs a message at ErrorLevel. The message includes any fields passed
func (f OverrideLogger) Errorf(format string, v ...interface{}) {}

// Sync flushes any buffered log entries.
func (f OverrideLogger) Sync() error {
	return nil
}

package logger

import (
	"log"
	"os"
)

// Logger interface for application logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

// StdLogger is a simple logger implementation using standard library
type StdLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

// New creates a new logger instance
func New() *StdLogger {
	return &StdLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an info message
func (l *StdLogger) Info(msg string, args ...interface{}) {
	l.infoLogger.Printf(msg, args...)
}

// Error logs an error message
func (l *StdLogger) Error(msg string, args ...interface{}) {
	l.errorLogger.Printf(msg, args...)
}

// Debug logs a debug message
func (l *StdLogger) Debug(msg string, args ...interface{}) {
	l.debugLogger.Printf(msg, args...)
}

// Warn logs a warning message
func (l *StdLogger) Warn(msg string, args ...interface{}) {
	l.warnLogger.Printf(msg, args...)
}

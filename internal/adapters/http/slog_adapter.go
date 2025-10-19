package http

import (
	"context"
	"log/slog"
)

// Logger defines the interface for logging
type Logger interface {
	Info(ctx context.Context, message string, fields map[string]interface{})
	Error(ctx context.Context, message string, fields map[string]interface{})
	Warn(ctx context.Context, message string, fields map[string]interface{})
	Debug(ctx context.Context, message string, fields map[string]interface{})
}

// SlogAdapter adapts slog.Logger to Logger interface
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlogAdapter creates a new adapter for slog
func NewSlogAdapter(logger *slog.Logger) *SlogAdapter {
	return &SlogAdapter{logger: logger}
}

// Info logs an informational message
func (l *SlogAdapter) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.logger.InfoContext(ctx, message, l.fieldsToAttrs(fields)...)
}

// Error logs an error
func (l *SlogAdapter) Error(ctx context.Context, message string, fields map[string]interface{}) {
	l.logger.ErrorContext(ctx, message, l.fieldsToAttrs(fields)...)
}

// Warn logs a warning
func (l *SlogAdapter) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	l.logger.WarnContext(ctx, message, l.fieldsToAttrs(fields)...)
}

// Debug logs a debug message
func (l *SlogAdapter) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	l.logger.DebugContext(ctx, message, l.fieldsToAttrs(fields)...)
}

// fieldsToAttrs converts map to slog.Attr
func (l *SlogAdapter) fieldsToAttrs(fields map[string]interface{}) []any {
	if len(fields) == 0 {
		return nil
	}

	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	return attrs
}

// NoOpLogger - no-op logger stub (for tests)
type NoOpLogger struct{}

func (n *NoOpLogger) Info(ctx context.Context, message string, fields map[string]interface{})  {}
func (n *NoOpLogger) Error(ctx context.Context, message string, fields map[string]interface{}) {}
func (n *NoOpLogger) Warn(ctx context.Context, message string, fields map[string]interface{})  {}
func (n *NoOpLogger) Debug(ctx context.Context, message string, fields map[string]interface{}) {}

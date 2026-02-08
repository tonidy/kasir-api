package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

// Logger wraps slog.Logger to provide a structured logging interface
type Logger struct {
	*slog.Logger
}

// New creates a new logger instance
func New(w io.Writer, level slog.Level) *Logger {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			// Customize timestamp format
			if attr.Key == "time" {
				return slog.Attr{
					Key:   attr.Key,
					Value: slog.StringValue(time.Now().Format("2006-01-02 15:04:05")),
				}
			}
			return attr
		},
	}

	logger := slog.New(slog.NewJSONHandler(w, opts))
	return &Logger{Logger: logger}
}

// NewDefault creates a new logger with default configuration
func NewDefault() *Logger {
	return New(os.Stdout, slog.LevelInfo)
}

// NewDev creates a new logger with development-friendly configuration
func NewDev() *Logger {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	return &Logger{Logger: logger}
}

// With adds attributes to the logger
func (l *Logger) With(attrs ...any) *Logger {
	return &Logger{Logger: l.Logger.With(attrs...)}
}

// WithContext adds context attributes to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	requestID := ctx.Value("request_id")
	if requestID != nil {
		return l.With("request_id", requestID)
	}
	return l
}

// Info logs an info message
func (l *Logger) Info(msg string, attrs ...any) {
	l.Logger.Info(msg, attrs...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, attrs ...any) {
	l.Logger.Warn(msg, attrs...)
}

// Error logs an error message
func (l *Logger) Error(msg string, attrs ...any) {
	l.Logger.Error(msg, attrs...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, attrs ...any) {
	l.Logger.Debug(msg, attrs...)
}

// InfoCtx logs an info message with context
func (l *Logger) InfoCtx(ctx context.Context, msg string, attrs ...any) {
	l.WithContext(ctx).Info(msg, attrs...)
}

// WarnCtx logs a warning message with context
func (l *Logger) WarnCtx(ctx context.Context, msg string, attrs ...any) {
	l.WithContext(ctx).Warn(msg, attrs...)
}

// ErrorCtx logs an error message with context
func (l *Logger) ErrorCtx(ctx context.Context, msg string, attrs ...any) {
	l.WithContext(ctx).Error(msg, attrs...)
}

// DebugCtx logs a debug message with context
func (l *Logger) DebugCtx(ctx context.Context, msg string, attrs ...any) {
	l.WithContext(ctx).Debug(msg, attrs...)
}

// LogRequest logs incoming HTTP requests
func (l *Logger) LogRequest(method, url, userAgent string, ip string, duration time.Duration) {
	l.Info("HTTP Request",
		slog.String("method", method),
		slog.String("url", url),
		slog.String("user_agent", userAgent),
		slog.String("ip", ip),
		slog.Duration("duration", duration),
	)
}

// LogError logs application errors with context
func (l *Logger) LogError(operation, err string, attrs ...any) {
	l.Error("Application Error",
		append(attrs,
			slog.String("operation", operation),
			slog.String("error", err),
		)...,
	)
}

// Global logger instance
var globalLogger *Logger

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		globalLogger = NewDefault()
	}
	return globalLogger
}

// Convenience functions using global logger
func Info(msg string, attrs ...any) {
	GetGlobalLogger().Info(msg, attrs...)
}

func Warn(msg string, attrs ...any) {
	GetGlobalLogger().Warn(msg, attrs...)
}

func Error(msg string, attrs ...any) {
	GetGlobalLogger().Error(msg, attrs...)
}

func Debug(msg string, attrs ...any) {
	GetGlobalLogger().Debug(msg, attrs...)
}

func InfoCtx(ctx context.Context, msg string, attrs ...any) {
	GetGlobalLogger().InfoCtx(ctx, msg, attrs...)
}

func WarnCtx(ctx context.Context, msg string, attrs ...any) {
	GetGlobalLogger().WarnCtx(ctx, msg, attrs...)
}

func ErrorCtx(ctx context.Context, msg string, attrs ...any) {
	GetGlobalLogger().ErrorCtx(ctx, msg, attrs...)
}

func DebugCtx(ctx context.Context, msg string, attrs ...any) {
	GetGlobalLogger().DebugCtx(ctx, msg, attrs...)
}

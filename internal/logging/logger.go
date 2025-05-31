package logging

import (
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger

// InitLogger initializes the global logger with the specified level and format
func InitLogger(logLevel string) {
	level := parseLogLevel(logLevel)
	
	// Use JSON format for production, text for development
	var handler slog.Handler
	if isProduction() {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}
	
	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// parseLogLevel converts string log level to slog.Level
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// isProduction determines if we're running in production based on environment
func isProduction() bool {
	env := os.Getenv("ENV")
	return env == "production" || env == "prod"
}

// Helper functions for common logging patterns
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

// WithFields returns a logger with the given fields pre-populated
func WithFields(args ...any) *slog.Logger {
	return Logger.With(args...)
}

// HTTP middleware logging helpers
func LogHTTPRequest(method, path string, statusCode int, args ...any) {
	baseArgs := []any{
		"method", method,
		"path", path,
		"status_code", statusCode,
	}
	allArgs := append(baseArgs, args...)
	
	if statusCode >= 400 {
		Logger.Warn("HTTP request completed with error", allArgs...)
	} else {
		Logger.Info("HTTP request completed", allArgs...)
	}
}
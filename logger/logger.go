package logger

import (
	"context"
	"log/slog"
)

type TraceIdKey struct{}

func InfoLog(ctx context.Context, message string) {
	if ctx != nil {
		slog.Info(ctx.Value(TraceIdKey{}).(string) + " " + message)
	} else {
		slog.Info(message)
	}
}

func WarningLog(ctx context.Context, message string) {
	if ctx != nil {
		slog.Warn(ctx.Value(TraceIdKey{}).(string) + " " + message)
	} else {
		slog.Warn(message)
	}
}

func ErrorLog(ctx context.Context, message string) {
	if ctx != nil {
		slog.Error(ctx.Value(TraceIdKey{}).(string) + " " + message)
	} else {
		slog.Error(message)
	}
}

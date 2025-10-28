package logger

import (
	"context"
	"log/slog"
)

type TraceIdKey struct{}

func InfoLog(ctx context.Context, message string) {
	slog.Info(ctx.Value(TraceIdKey{}).(string) + " " + message)
}

func WarningLog(ctx context.Context, message string) {
	slog.Warn(ctx.Value(TraceIdKey{}).(string) + " " + message)
}

func ErrorLog(ctx context.Context, message string) {
	slog.Error(ctx.Value(TraceIdKey{}).(string) + " " + message)
}

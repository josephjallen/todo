package logger

import (
	"context"
	"log/slog"
)

type TraceIdKey struct{}

func GetCtxLogger(ctx context.Context) *slog.Logger {
	traceID := ctx.Value(TraceIdKey{})
	if traceID == nil {
		return slog.Default()
	}
	return slog.With("trace_id", traceID.(string))
}

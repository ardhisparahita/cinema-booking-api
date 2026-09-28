package logger

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

func New() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	return slog.New(handler)
}

func FromContext(ctx context.Context, base *slog.Logger) *slog.Logger {
	spanContext := trace.SpanContextFromContext(ctx)

	if !spanContext.IsValid() {
		return base
	}

	return base.With(
		"trace_id", spanContext.TraceID().String(),
		"span_id", spanContext.SpanID().String(),
	)
}

package middleware_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	appMiddleware "github.com/ardhisparahita/cinema-booking-api/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestRequestLogger_TraceIDMatchesSpan(t *testing.T) {

	exporter := tracetest.NewInMemoryExporter()

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(
			sdktrace.AlwaysSample(),
		),
		sdktrace.WithSpanProcessor(
			sdktrace.NewSimpleSpanProcessor(exporter),
		),
	)

	defer func() {
		_ = tracerProvider.Shutdown(context.Background())
	}()

	tracer := tracerProvider.Tracer("request-logger-test")

	var logBuffer bytes.Buffer

	handler := slog.NewTextHandler(
		&logBuffer,
		&slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: false,
		},
	)

	logger := slog.New(handler)

	app := fiber.New()

	app.Use(requestid.New())

	app.Use(func(c fiber.Ctx) error {

		ctx, span := tracer.Start(
			c.Context(),
			"HTTP GET /test",
		)

		c.SetContext(ctx)

		defer span.End()

		return c.Next()
	})

	app.Use(
		appMiddleware.RequestLogger(logger),
	)

	app.Get("/test", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(
		"GET",
		"/test",
		nil,
	)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			resp.StatusCode,
		)
	}

	spans := exporter.GetSpans()

	if len(spans) != 1 {
		t.Fatalf(
			"expected 1 span, got %d",
			len(spans),
		)
	}

	traceID := spans[0].SpanContext.TraceID().String()

	if traceID == "" {
		t.Fatal("trace ID is empty")
	}

	logOutput := logBuffer.String()

	if logOutput == "" {
		t.Fatal("expected log output, got empty string")
	}

	expectedTraceID := "trace_id=" + traceID

	if !strings.Contains(logOutput, expectedTraceID) {
		t.Fatalf(
			"trace ID not found in log\nexpected: %s\nlog: %s",
			expectedTraceID,
			logOutput,
		)
	}

	t.Logf("Trace ID: %s", traceID)
	t.Logf("Log: %s", logOutput)
}

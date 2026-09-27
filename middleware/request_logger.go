package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func RequestLogger(log *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		status := c.Response().StatusCode()

		attrs := []any{
			"request_id", requestid.FromContext(c),
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"latency", time.Since(start),
		}

		if userID, ok := c.Locals("user_id").(uint); ok {
			attrs = append(attrs, "user_id", userID)
		}

		switch {
		case status >= 500:
			if err != nil {
				attrs = append(attrs, "error", err)
			}
			log.Error("request completed with server error", attrs...)

		case status >= 400:
			log.Error("request completed with client error", attrs...)

		default:
			log.Info("request completed", attrs...)
		}

		return err
	}
}

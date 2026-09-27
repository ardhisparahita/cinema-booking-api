package utils

import (
	"errors"
	"log/slog"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"

	var fiberErr *fiber.Error

	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	}

	requestID := requestid.FromContext(c)

	attrs := []any{
		"request_id", requestID,
		"method", c.Method(),
		"path", c.Path(),
		"status", code,
		"error", err,
	}

	if code >= 500 {
		slog.Error("unhandled request error", attrs...)
	} else {
		slog.Warn("request error", attrs...)
	}

	return c.Status(code).JSON(response.ErrorResponse{
		Code:    code,
		Status:  "error",
		Message: message,
	})
}

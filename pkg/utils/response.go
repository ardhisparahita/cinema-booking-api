package utils

import (
	"errors"
	"log/slog"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	applogger "github.com/ardhisparahita/cinema-booking-api/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func ResponseSuccess(c fiber.Ctx, code int, message string, data any) error {
	return c.Status(code).JSON(response.WebResponse{
		Code:    code,
		Status:  "Success",
		Message: message,
		Data:    data,
	})
}

func ResponseError(c fiber.Ctx, code int, message string, err error) error {
	var validationErrs validator.ValidationErrors

	if errors.As(err, &validationErrs) {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  "Error",
			Message: "validation failed",
			Data:    ValidationError(validationErrs),
		})
	}

	logWithTrace := applogger.FromContext(c.Context(), slog.Default())

	loggerAttrs := []any{
		"request_id", requestid.FromContext(c),
		"method", c.Method(),
		"path", c.Path(),
		"status", code,
		"message", message,
	}

	if err != nil {
		loggerAttrs = append(loggerAttrs, "error", err)
	}

	switch {
	case code >= 500:
		logWithTrace.Error("api error", loggerAttrs...)
	case code >= 400:
		logWithTrace.Warn("api error", loggerAttrs...)
	default:
		logWithTrace.Info("api error", loggerAttrs...)
	}

	return c.Status(code).JSON(response.WebResponse{
		Code:    code,
		Status:  "Error",
		Message: message,
		Data:    nil,
	})
}

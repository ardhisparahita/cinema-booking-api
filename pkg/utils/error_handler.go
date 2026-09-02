package utils

import (
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(response.ErrorResponse{
		Code:    code,
		Status:  "error",
		Message: err.Error(),
	})
}

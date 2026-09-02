package utils

import (
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
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

	errData := ""
	if err != nil {
		errData = err.Error()
	}

	return c.Status(code).JSON(response.WebResponse{
		Code:    code,
		Status:  "Error",
		Message: message,
		Data:    errData,
	})
}

package handler

import (
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/middleware"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) GetProfile(c fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey)

	id, ok := userID.(uint)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user identity",
		)
	}

	res, err := h.UserService.GetProfile(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"user not found",
				err,
			)
		}

		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"internal server error",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"profile retrieved successfully",
		res,
	)
}

func (h *UserHandler) UpdateProfile(c fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey)

	id, ok := userID.(uint)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user identify",
		)
	}

	var req request.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid request body",
			err,
		)
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"validation failed",
			err,
		)
	}

	res, err := h.UserService.UpdateProfile(
		c.Context(),
		id,
		req,
	)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"user not found",
				err,
			)
		}
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"internal server error",
			err,
		)
	}
	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"profile updated successfully",
		res,
	)
}

func (h *UserHandler) DeleteProfile(c fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey)

	id, ok := userID.(uint)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user identify",
		)
	}

	if err := h.UserService.DeleteProfile(
		c.Context(),
		id,
	); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"user not found",
				err,
			)
		}

		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"internal server error",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"account delete successfully",
		nil,
	)
}

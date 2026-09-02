package handler

import (
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req request.RegisterRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "invalid request body", err)
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "validation failed", err)
	}

	res, err := h.authService.Register(c.Context(), req)
	if err != nil {
		if err.Error() == "email already exists" {
			return utils.ResponseError(c, fiber.StatusConflict, "registration failed", err)
		}
		return utils.ResponseError(c, fiber.StatusInternalServerError, "internal server error", err)
	}

	return utils.ResponseSuccess(c, fiber.StatusCreated, "user registered successfully", res)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req request.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "invalid request body", err)
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "validation failed", err)
	}

	res, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusUnauthorized, "login failed", err)
	}

	return utils.ResponseSuccess(c, fiber.StatusOK, "login successful", res)
}

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req request.RefreshToken

	if err := c.Bind().Body(&req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "invalid request body", err)
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "validation failed", err)
	}

	res, err := h.authService.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusUnauthorized, "token refresh failed", err)
	}

	return utils.ResponseSuccess(c, fiber.StatusOK, "token refreshed successfully", res)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var req request.RefreshToken

	if err := c.Bind().Body(&req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "invalid request body", err)
	}

	if err := utils.ValidationStruct(req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "validation failed", err)
	}

	if err := h.authService.Logout(c.Context(), req.RefreshToken); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "logout failed", err)
	}
	return utils.ResponseSuccess(c, fiber.StatusOK, "logout successful", nil)
}

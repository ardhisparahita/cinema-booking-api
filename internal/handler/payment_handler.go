package handler

import (
	"errors"
	"strconv"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	appErrors "github.com/ardhisparahita/cinema-booking-api/internal/errors"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/middleware"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type PaymentHandler struct {
	Service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		Service: service,
	}
}

func (h *PaymentHandler) Create(c fiber.Ctx) error {
	var req request.CreatePaymentRequest

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

	userID, ok := c.Locals(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return utils.ResponseError(
			c,
			fiber.StatusUnauthorized,
			"user not authenticated",
			nil,
		)
	}

	res, err := h.Service.CreatePayment(c.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrBookingNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"booking not found",
				err,
			)
		case errors.Is(err, appErrors.ErrBookingNotPending):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"booking is not pending",
				err,
			)
		case errors.Is(err, appErrors.ErrBookingExpired):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"booking has expired",
				err,
			)
		case errors.Is(err, appErrors.ErrPaymentAlreadyExist):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"payment already exist",
				err,
			)
		case errors.Is(err, appErrors.ErrInvalidPaymentMethod):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"invalid payment method",
				err,
			)

		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to create payment",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"payment created successfully",
		res,
	)
}

func (h *PaymentHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid payment id",
			err,
		)
	}

	userID, ok := c.Locals(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return utils.ResponseError(
			c,
			fiber.StatusUnauthorized,
			"user not authenticated",
			nil,
		)
	}

	res, err := h.Service.GetPaymentByID(c.Context(), userID, uint(id))
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrPaymentNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"payment not found",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to get payment",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"payment retrieved successfully",
		res,
	)
}

func (h *PaymentHandler) Confirm(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid payment id",
			err,
		)
	}

	userID, ok := c.Locals(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return utils.ResponseError(
			c,
			fiber.StatusUnauthorized,
			"user not authenticated",
			nil,
		)
	}

	res, err := h.Service.ConfirmPayment(c.Context(), userID, uint(id))
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrPaymentNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"payment not found",
				err,
			)
		case errors.Is(err, appErrors.ErrPaymentAlreadyPaid):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"payment already paid",
				err,
			)
		case errors.Is(err, appErrors.ErrBookingNotPending):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"booking is not pending",
				err,
			)
		case errors.Is(err, appErrors.ErrBookingExpired):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"booking has expired",
				err,
			)
		case errors.Is(err, appErrors.ErrPaymentCannotConfirm):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"payment cannot be confirmed",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to confirm payment",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"payment confirmed successfully",
		res,
	)
}

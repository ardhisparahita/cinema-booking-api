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

type BookingHandler struct {
	Service service.BookingService
}

func NewBookingHandler(service service.BookingService) *BookingHandler {
	return &BookingHandler{
		Service: service,
	}
}

func (h *BookingHandler) Create(c fiber.Ctx) error {
	var req request.CreateBookingRequest

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

	res, err := h.Service.CreateBooking(c.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidBookingSeats):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"at least one seat in required",
				err,
			)
		case errors.Is(err, appErrors.ErrDuplicateSeat):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"duplicate seat in booking",
				err,
			)
		case errors.Is(err, appErrors.errshow):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"showtime not found",
				err,
			)
		case errors.Is(err, appErrors.ErrShowtimeFinished):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"showtime already finished",
				err,
			)
		case errors.Is(err, appErrors.ErrSeatNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"one or more seats not found",
				err,
			)
		case errors.Is(err, appErrors.ErrSeatWrongStudio):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"one or more seats do not belong to showtime studio",
				err,
			)
		case errors.Is(err, appErrors.ErrSeatAlreadyBooked):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"one or more seats already booked",
				err,
			)
		case errors.Is(err, appErrors.ErrSeatLocked):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"one or more seats are currently locked",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to create booking",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"booking created successfully",
		res,
	)
}

func (h *BookingHandler) GetMyBookings(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return utils.ResponseError(
			c,
			fiber.StatusUnauthorized,
			"user not authenticated",
			nil,
		)
	}
	res, err := h.Service.GetMyBookings(c.Context(), userID)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get bookings",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"bookings retrieved successfully",
		res,
	)
}

func (h *BookingHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid booking id",
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

	res, err := h.Service.GetBookingByID(c.Context(), userID, uint(id))
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrBookingNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"booking not found",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to get booking",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"booking retrieved successfully",
		res,
	)
}

func (h *BookingHandler) CancelBooking(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid booking id",
			err,
		)
	}

	userID, ok := c.Locals(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return utils.ResponseError(
			c,
			fiber.StatusUnauthorized,
			"user not authenticated",
			err,
		)
	}

	err = h.Service.CancelBooking(c.Context(), userID, uint(id))
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrBookingNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"booking not found",
				err,
			)
		case errors.Is(err, appErrors.ErrBookingCannotCancel):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"booking cannot be cancelled",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to cancel booking",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"booking cancelled successfully",
		nil,
	)
}

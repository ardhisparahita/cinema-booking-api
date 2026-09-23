package handler

import (
	"errors"
	"strconv"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	appErrors "github.com/ardhisparahita/cinema-booking-api/internal/errors"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type SeatHandler struct {
	Service service.SeatService
}

func NewSeatHandler(service service.SeatService) *SeatHandler {
	return &SeatHandler{
		Service: service,
	}
}

func (h *SeatHandler) GetAll(c fiber.Ctx) error {
	studioID, err := strconv.ParseUint(c.Params("studioId"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid studio id",
			err,
		)
	}

	res, err := h.Service.GetAllSeats(c.Context(), uint(studioID))

	if err != nil {
		if errors.Is(err, appErrors.ErrStudioNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"studio not found",
				err,
			)
		}

		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get seats",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"seats retrieved successfully",
		res,
	)
}

func (h *SeatHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid seat id",
			err,
		)
	}

	res, err := h.Service.GetSetByID(c.Context(), uint(id))
	if err != nil {
		if errors.Is(err, appErrors.ErrSeatNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"seat not found",
				err,
			)
		}
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get seat",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"seat retrieved successfully",
		res,
	)
}

func (h *SeatHandler) Create(c fiber.Ctx) error {
	studioID, err := strconv.ParseUint(c.Params("studioId"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid studio id",
			err,
		)
	}

	var req request.CreateAndUpdateSeatRequest

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

	res, err := h.Service.CreateSeat(c.Context(), uint(studioID), req)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrStudioNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"studio not found",
				err,
			)
		case errors.Is(err, appErrors.ErrSeatAlreadyExists):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"seat already exist",
				err,
			)
		case errors.Is(err, appErrors.ErrInvalidSeatType), errors.Is(err, appErrors.ErrInvalidSeatRow):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"invalid seat data",
				err,
			)
		case errors.Is(err, appErrors.ErrSeatOutsideStudioLayout):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"seat position is outside studio layout",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to create seat",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"seat created successfully",
		res,
	)
}

func (h *SeatHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid seat id",
			err,
		)
	}

	var req request.CreateAndUpdateSeatRequest

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

	res, err := h.Service.UpdateSeat(c.Context(), uint(id), req)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrSeatNotFound):
			utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"seat not found",
				err,
			)
		case errors.Is(err, appErrors.ErrSeatAlreadyExists):
			utils.ResponseError(
				c,
				fiber.StatusConflict,
				"seat already exist",
				err,
			)
		case errors.Is(err, appErrors.ErrInvalidSeatType), errors.Is(err, appErrors.ErrInvalidSeatRow):
			utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"invalid seat data",
				err,
			)
		case errors.Is(err, appErrors.ErrSeatOutsideStudioLayout):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"seat position is outside studio layout",
				err,
			)
		default:
			utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to update seat",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"seat updated successfully",
		res,
	)
}

func (h *SeatHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusNotFound,
			"invalid seat id",
			err,
		)
	}

	if err := h.Service.DeleteSeat(c.Context(), uint(id)); err != nil {
		if errors.Is(err, appErrors.ErrSeatNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"seat not found",
				err,
			)
		}
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to delete seat",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"seat deleted successfully",
		nil,
	)
}

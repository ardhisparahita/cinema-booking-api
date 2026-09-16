package handler

import (
	"errors"
	"strconv"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type ShowtimeHandler struct {
	Service service.ShowtimeService
}

func NewShowtimeHandler(service service.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{
		Service: service,
	}
}

func (h *ShowtimeHandler) GetAll(c fiber.Ctx) error {
	res, err := h.Service.GetAllShowtimes(c.Context())
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get showtimes",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"showtimes retrieved successfully",
		res,
	)
}

func (h *ShowtimeHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid showtime id",
			err,
		)
	}

	res, err := h.Service.GetShowtimeByID(c.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrShowtimeNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"showtime not found",
				err,
			)
		}
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get showtime",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"showtime retrieved successfully",
		res,
	)
}

func (h *ShowtimeHandler) Create(c fiber.Ctx) error {
	var req request.CreateAndUpdateShowtimeRequest

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

	res, err := h.Service.CreateShowtime(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMovieNotFoundShowtime):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"movie not found",
				err,
			)
		case errors.Is(err, service.ErrStudioNotFoundShowtime):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"studio not found",
				err,
			)
		case errors.Is(err, service.ErrInvalidShowtimeTime), errors.Is(err, service.ErrInvalidShowtimePrice):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"invalid showtime data",
				err,
			)
		case errors.Is(err, service.ErrShowtimeConflict):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"showtime conflicts with another showtime",
				err,
			)
		case errors.Is(err, service.ErrShowtimeInPast):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"showtime cannot be in the past",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to create showtime",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"showtime created successfully",
		res,
	)
}

func (h *ShowtimeHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid showtime id",
			err,
		)
	}

	var req request.CreateAndUpdateShowtimeRequest

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

	res, err := h.Service.UpdateShowtime(c.Context(), uint(id), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrShowtimeNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"showtime not found",
				err,
			)
		case errors.Is(err, service.ErrMovieNotFoundShowtime):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"movie not found",
				err,
			)
		case errors.Is(err, service.ErrStudioNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"studio not found",
				err,
			)
		case errors.Is(err, service.ErrInvalidShowtimeTime), errors.Is(err, service.ErrInvalidShowtimePrice):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"invalid showtime data",
				err,
			)
		case errors.Is(err, service.ErrShowtimeConflict):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"showtime conflicts with another showtime",
				err,
			)
		case errors.Is(err, service.ErrShowtimeInPast):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"showtime cannot be in the past",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to update showtime",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"showtime updated successfully",
		res,
	)
}

func (h *ShowtimeHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid showtime id",
			err,
		)
	}

	if err := h.Service.DeleteShowtime(c.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrShowtimeNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"showtime not found",
				err,
			)
		}
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to delete showtime",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"showtime deleted successfully",
		nil,
	)
}

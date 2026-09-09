package handler

import (
	"errors"
	"strconv"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type MovieHandler struct {
	MovieService service.MovieService
}

func NewMovieHandler(movieService service.MovieService) *MovieHandler {
	return &MovieHandler{
		MovieService: movieService,
	}
}

func (h *MovieHandler) GetAll(c fiber.Ctx) error {
	res, err := h.MovieService.GetAllMovies(c.Context())
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get movies",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"movie retrieved successfully",
		res,
	)
}

func (h *MovieHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid movie id",
			err,
		)
	}

	res, err := h.MovieService.GetMovieByID(c.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"movie not found",
				err,
			)
		}

		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get movie",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"movie retrieved successfully",
		res,
	)
}

func (h *MovieHandler) Create(c fiber.Ctx) error {
	var req request.CreateMovieRequest

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
			"invalid request body",
			err,
		)
	}

	res, err := h.MovieService.CreateMovie(c.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrGenreNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"one or more genres not found",
				err,
			)
		case errors.Is(err, service.ErrInvalidMovieRating):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"invalid movie rating",
				err,
			)

		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to create movie",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"movie created successfully",
		res,
	)
}

func (h *MovieHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid movie id",
			err,
		)
	}

	var req request.UpdateMovieRequest

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
			"invalid request body",
			err,
		)
	}

	res, err := h.MovieService.UpdateMovie(c.Context(), uint(id), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMovieNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"movie not found",
				err,
			)

		case errors.Is(err, service.ErrGenreNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"one or more genres not found",
				err,
			)
		case errors.Is(err, service.ErrInvalidMovieRating):
			return utils.ResponseError(
				c,
				fiber.StatusBadRequest,
				"invalid movie rating",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to update movie",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"movie updated successfully",
		res,
	)
}

func (h *MovieHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid movie id",
			err,
		)
	}

	if err := h.MovieService.DeleteMovie(c.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"movie not found",
				err,
			)
		}

		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to delete movie",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"movie deleted successfully",
		nil,
	)
}

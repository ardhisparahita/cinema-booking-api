package handler

import (
	"errors"
	"strconv"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type GenreHandler struct {
	GenreService service.GenreService
}

func NewGenreHandler(genreService service.GenreService) *GenreHandler {
	return &GenreHandler{
		GenreService: genreService,
	}
}

func (h *GenreHandler) GetAll(c fiber.Ctx) error {
	res, err := h.GenreService.GetAllGenres(c.Context())

	if err != nil {
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
		"genre retrieved successfully",
		res,
	)
}

func (h *GenreHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseInt(
		c.Params("id"),
		10,
		64,
	)

	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid genre id",
			err,
		)
	}

	res, err := h.GenreService.GetGenreByID(c.Context(), uint(id))

	if err != nil {
		if errors.Is(err, service.ErrGenreNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"genre not found",
				err,
			)
		}

		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get genre",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"genre retrieved successfully",
		res,
	)
}

func (h *GenreHandler) Create(c fiber.Ctx) error {
	var req request.CreateGenreRequest

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

	res, err := h.GenreService.CreateGenre(c.Context(), req)

	if err != nil {
		if errors.Is(err, service.ErrGenreAlreadyExist) {
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"genre already exist",
				err,
			)
		}

		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to create genre",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"genre created successfully",
		res,
	)
}

func (h *GenreHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)

	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid genre id",
			err,
		)
	}

	var req request.UpdateGenreRequest

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

	res, err := h.GenreService.UpdateGenre(c.Context(), uint(id), req)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrGenreNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"genre not found",
				err,
			)
		case errors.Is(err, service.ErrGenreAlreadyExist):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"genre already exist",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to update genre",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"genre update successfully",
		res,
	)
}

func (h *GenreHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid genre id",
			err,
		)
	}

	if err := h.GenreService.DeleteGenre(c.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrGenreNotFound) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"genre not found",
				err,
			)
		}

		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to delete genre",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"genre delete successfully",
		nil,
	)
}

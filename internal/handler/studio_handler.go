package handler

import (
	"errors"
	"strconv"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type StudioHandler struct {
	Service service.StudioService
}

func NewStudioHandler(service service.StudioService) *StudioHandler {
	return &StudioHandler{
		Service: service,
	}
}

func (h *StudioHandler) GetAll(c fiber.Ctx) error {
	theaterID, err := strconv.ParseUint(c.Params("theaterId"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid theater id",
			err,
		)
	}

	res, err := h.Service.GetAllStudios(c.Context(), uint(theaterID))

	if err != nil {
		if errors.Is(err, service.ErrTheaterNotFoundForStudio) {
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"theater not found",
				err,
			)
		}
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get studios",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"studios retrieved successfully",
		res,
	)
}

func (h *StudioHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid studio id",
			err,
		)
	}

	res, err := h.Service.GetStudioByID(c.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrStudioNotFound) {
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
			"failed to get studio",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"studio retrieved successfully",
		res,
	)
}

func (h *StudioHandler) Create(c fiber.Ctx) error {
	theaterID, err := strconv.ParseUint(c.Params("theaterId"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid theater id",
			err,
		)
	}

	var req request.CreateAndUpdateStudioRequest

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

	res, err := h.Service.CreateStudio(c.Context(), uint(theaterID), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTheaterNotFoundForStudio):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"theater not found",
				err,
			)
		case errors.Is(err, service.ErrStudioAlreadyExist):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"studio already exist",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to create studio",
				err,
			)
		}
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"studio created successfully",
		res,
	)
}

func (h *StudioHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid studio id",
			err,
		)
	}

	var req request.CreateAndUpdateStudioRequest

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

	res, err := h.Service.UpdateStudio(c.Context(), uint(id), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStudioNotFound):
			return utils.ResponseError(
				c,
				fiber.StatusNotFound,
				"studio not found",
				err,
			)
		case errors.Is(err, service.ErrStudioAlreadyExist):
			return utils.ResponseError(
				c,
				fiber.StatusConflict,
				"studio already exist",
				err,
			)
		default:
			return utils.ResponseError(
				c,
				fiber.StatusInternalServerError,
				"failed to update studio",
				err,
			)
		}
	}
	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"studio updated successfully",
		res,
	)
}

func (h *StudioHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid studio id",
			err,
		)
	}

	if err := h.Service.DeleteStudio(c.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrStudioNotFound) {
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
			"failed to delete studio",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"studio deleted successfully",
		nil,
	)
}

package handler

import (
	"errors"
	"strconv"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	appErrors "github.com/ardhisparahita/cinema-booking-api/internal/errors"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type TheaterHandler struct {
	Service service.TheaterService
}

func NewTheaterHandler(service service.TheaterService) *TheaterHandler {
	return &TheaterHandler{
		Service: service,
	}
}

func (h *TheaterHandler) GetAll(c fiber.Ctx) error {
	var query request.PaginationRequest

	if err := c.Bind().Query(&query); err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid pagination parameter",
			err,
		)
	}

	pagination := utils.GetPagination(query.Page, query.Limit)

	res, total, err := h.Service.GetAllTheaters(c.Context(), pagination.Page, pagination.Limit)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to get theaters",
			err,
		)
	}

	data := struct {
		Data       []response.TheaterResponse  `json:"data"`
		Pagination response.PaginationResponse `json:"pagination"`
	}{
		Data: res,
		Pagination: response.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.GetTotalPage(total, pagination.Limit),
		},
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"theaters retrieved successfully",
		data,
	)
}

func (h *TheaterHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid theater id",
			err,
		)
	}

	res, err := h.Service.GetTheaterByID(c.Context(), uint(id))

	if err != nil {
		if errors.Is(err, appErrors.ErrTheaterNotFound) {
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
			"failed to get theater",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"theater retrieved successfully",
		res,
	)
}

func (h *TheaterHandler) Create(c fiber.Ctx) error {
	var req request.CreateAndUpdateTheaterRequest

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

	res, err := h.Service.CreateTheater(c.Context(), req)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusInternalServerError,
			"failed to create theater",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusCreated,
		"theater created successfully",
		res,
	)
}

func (h *TheaterHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid theater id",
			err,
		)
	}

	var req request.CreateAndUpdateTheaterRequest

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

	res, err := h.Service.UpdateTheater(c.Context(), uint(id), req)
	if err != nil {
		if errors.Is(err, appErrors.ErrTheaterNotFound) {
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
			"failed to update theater",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"theater updated successfully",
		res,
	)
}

func (h *TheaterHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ResponseError(
			c,
			fiber.StatusBadRequest,
			"invalid theater id",
			err,
		)
	}

	if err := h.Service.DeleteTheater(c.Context(), uint(id)); err != nil {
		if errors.Is(err, appErrors.ErrTheaterNotFound) {
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
			"failed to delete theater",
			err,
		)
	}

	return utils.ResponseSuccess(
		c,
		fiber.StatusOK,
		"theater deleted successfully",
		nil,
	)
}

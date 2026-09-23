package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	appErrors "github.com/ardhisparahita/cinema-booking-api/internal/errors"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
)

// var (
// 	ErrSeatNotFound            = errors.New("seat not found")
// 	ErrSeatAlreadyExist        = errors.New("seat already exist")
// 	ErrStudioNotFoundSeat      = errors.New("studio not found")
// 	ErrInvalidSeatType         = errors.New("invalid seat type")
// 	ErrInvalidSeatRow          = errors.New("invalid seat row")
// 	ErrInvalidSeatColumn       = errors.New("invalid seat column")
// 	ErrSeatOutsideStudioLayout = errors.New("seat position is outside studio layout")
// )

type SeatServiceImpl struct {
	Repo       repository.SeatRepository
	StudioRepo repository.StudioRepository
}

func NewSeatService(repo repository.SeatRepository, studioRepo repository.StudioRepository) SeatService {
	return &SeatServiceImpl{
		Repo:       repo,
		StudioRepo: studioRepo,
	}
}

func (s *SeatServiceImpl) CreateSeat(ctx context.Context, studioID uint, req request.CreateAndUpdateSeatRequest) (*response.SeatResponse, error) {
	studio, err := s.StudioRepo.FindStudioByID(ctx, studioID)
	if err != nil {
		if errors.Is(err, appErrors.ErrStudioNotFound) {
			return nil, appErrors.ErrStudioNotFound
		}

		return nil, err
	}

	rowLabel := strings.ToUpper(strings.TrimSpace(req.RowLabel))
	if rowLabel == "" || len(rowLabel) > 2 {
		return nil, appErrors.ErrInvalidSeatRow
	}

	if req.SeatType != "regular" && req.SeatType != "vip" {
		return nil, appErrors.ErrInvalidSeatType
	}

	rowNumber, err := rowLabelToNumber(rowLabel)
	if err != nil {
		return nil, appErrors.ErrInvalidSeatRow
	}

	if uint16(rowNumber) > studio.TotalRows {
		return nil, appErrors.ErrSeatOutsideStudioLayout
	}

	if req.ColNumber > studio.TotalCols {
		return nil, appErrors.ErrSeatOutsideStudioLayout
	}

	existing, err := s.Repo.FindSeatByPosition(ctx, studioID, rowLabel, req.ColNumber)
	if err == nil && existing != nil {
		return nil, appErrors.ErrSeatAlreadyExists
	}
	if err != nil && !errors.Is(err, appErrors.ErrSeatNotFound) {
		return nil, err
	}

	seat := &models.Seat{
		StudioID:  studioID,
		RowLabel:  rowLabel,
		ColNumber: req.ColNumber,
		SeatType:  req.SeatType,
	}

	if err := s.Repo.CreateSeat(ctx, seat); err != nil {
		if errors.Is(err, appErrors.ErrSeatAlreadyExists) {
			return nil, appErrors.ErrSeatAlreadyExists
		}
		return nil, err
	}

	return toSeatResponse(seat), nil
}

func (s *SeatServiceImpl) GetAllSeats(ctx context.Context, studioID uint) ([]response.SeatResponse, error) {
	if _, err := s.StudioRepo.FindStudioByID(ctx, studioID); err != nil {
		if errors.Is(err, appErrors.ErrStudioNotFound) {
			return nil, appErrors.ErrStudioNotFound
		}
		return nil, err
	}

	seats, err := s.Repo.FindAllSeats(ctx, studioID)
	if err != nil {
		return nil, err
	}

	result := make([]response.SeatResponse, 0, len(seats))

	for _, seat := range seats {
		result = append(result, *toSeatResponse(&seat))
	}

	return result, nil
}

func (s *SeatServiceImpl) GetSetByID(ctx context.Context, id uint) (*response.SeatResponse, error) {
	seat, err := s.Repo.FindSeatByID(ctx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrSeatNotFound) {
			return nil, appErrors.ErrSeatNotFound
		}
		return nil, err
	}

	return toSeatResponse(seat), nil
}

func (s *SeatServiceImpl) UpdateSeat(ctx context.Context, id uint, req request.CreateAndUpdateSeatRequest) (*response.SeatResponse, error) {
	seat, err := s.Repo.FindSeatByID(ctx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrSeatNotFound) {
			return nil, appErrors.ErrSeatNotFound
		}
		return nil, err
	}

	studio, err := s.StudioRepo.FindStudioByID(ctx, seat.StudioID)
	if err != nil {
		if errors.Is(err, appErrors.ErrStudioNotFound) {
			return nil, appErrors.ErrStudioNotFound
		}
		return nil, err
	}

	rowLabel := strings.ToUpper(strings.TrimSpace(req.RowLabel))

	if rowLabel == "" || len(rowLabel) > 2 {
		return nil, appErrors.ErrInvalidSeatRow
	}

	if req.SeatType != "regular" && req.SeatType != "vip" {
		return nil, appErrors.ErrInvalidSeatType
	}

	rowNumber, err := rowLabelToNumber(rowLabel)
	if err != nil {
		return nil, appErrors.ErrInvalidSeatRow
	}

	if uint16(rowNumber) > studio.TotalRows {
		return nil, appErrors.ErrSeatOutsideStudioLayout
	}
	if req.ColNumber > studio.TotalCols {
		return nil, appErrors.ErrSeatOutsideStudioLayout
	}

	existing, err := s.Repo.FindSeatByPosition(ctx, seat.StudioID, rowLabel, req.ColNumber)
	if err == nil && existing.ID != seat.ID {
		return nil, appErrors.ErrSeatAlreadyExists
	}

	if err != nil && !errors.Is(err, appErrors.ErrSeatNotFound) {
		return nil, err
	}

	seat.RowLabel = rowLabel
	seat.ColNumber = req.ColNumber
	seat.SeatType = req.SeatType

	if err := s.Repo.UpdateSeat(ctx, seat); err != nil {
		return nil, err
	}

	return toSeatResponse(seat), nil
}

func (s *SeatServiceImpl) DeleteSeat(ctx context.Context, id uint) error {
	if _, err := s.Repo.FindSeatByID(ctx, id); err != nil {
		if errors.Is(err, appErrors.ErrSeatNotFound) {
			return appErrors.ErrSeatNotFound
		}
		return err
	}

	return s.Repo.DeleteSeat(ctx, uint16(id))
}

func toSeatResponse(seat *models.Seat) *response.SeatResponse {
	return &response.SeatResponse{
		ID:        seat.ID,
		StudioID:  seat.StudioID,
		RowLabel:  seat.RowLabel,
		ColNumber: seat.ColNumber,
		SeatType:  seat.SeatType,
		CreatedAt: seat.CreatedAt,
	}
}

func rowLabelToNumber(rowLabel string) (int, error) {
	if len(rowLabel) == 0 || len(rowLabel) > 2 {
		return 0, errors.New("invalid row label")
	}

	result := 0

	for _, char := range rowLabel {
		if char < 'A' || char > 'Z' {
			return 0, errors.New("invalid row label")
		}

		result = result*26 + int(char-'A'+1)
	}

	return result, nil
}

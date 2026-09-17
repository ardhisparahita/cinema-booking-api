package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrBookingNotFound      = errors.New("booking not found")
	ErrShowtimeNotFoundBook = errors.New("showtime not found")
	ErrShowtimeFinished     = errors.New("showtime already finished")
	ErrSeatNotFoundBook     = errors.New("one or more seats not found")
	ErrSeatWrongStudio      = errors.New("one or more seats do not belong to showtime studio")
	ErrSeatAlreadyBooked    = errors.New("one or more seats already booked")
	ErrDuplicateSeat        = errors.New("duplicate seat in booking")
	ErrInvalidBookingSeats  = errors.New("at least on seat is required")
	ErrBookingCannotCancel  = errors.New("booking cannot be cancelled")
)

type BookingServiceImpl struct {
	Repo         repository.BookingRepository
	ShowtimeRepo repository.ShowtimeRepository
	SeatRepo     repository.SeatRepository
	DB           *gorm.DB
}

func NewBookingService(repo repository.BookingRepository, showtimeRepo repository.ShowtimeRepository, seatRepo repository.SeatRepository, db *gorm.DB) BookingService {
	return &BookingServiceImpl{
		Repo:         repo,
		ShowtimeRepo: showtimeRepo,
		SeatRepo:     seatRepo,
		DB:           db,
	}
}

func (s *BookingServiceImpl) CreateBooking(ctx context.Context, userID uint, req request.CreateBookingRequest) (*response.BookingResponse, error) {
	if len(req.SeatIDs) == 0 {
		return nil, ErrInvalidBookingSeats
	}

	seatIDs := uniqueUint(req.SeatIDs)
	if len(seatIDs) != len(req.SeatIDs) {
		return nil, ErrDuplicateSeat
	}

	showtime, err := s.ShowtimeRepo.FindShowtimeByID(ctx, req.ShowtimeID)
	if err != nil {
		if errors.Is(err, repository.ErrShowtimeNotFound) {
			return nil, ErrShowtimeNotFoundBook
		}
		return nil, err
	}

	if !showtime.EndTime.After(time.Now()) {
		return nil, ErrShowtimeFinished
	}

	seats, err := s.SeatRepo.FindSeatByIDs(ctx, seatIDs)
	if err != nil {
		return nil, err
	}

	if len(seats) != len(seatIDs) {
		return nil, ErrSeatNotFoundBook
	}

	for _, seat := range seats {
		if seat.StudioID != showtime.StudioID {
			return nil, ErrSeatWrongStudio
		}
	}

	bookedSeatIDs, err := s.Repo.FindBookedSeatIDs(ctx, req.ShowtimeID, seatIDs)
	if err != nil {
		return nil, err
	}

	if len(bookedSeatIDs) > 0 {
		return nil, ErrSeatAlreadyBooked
	}

	totalPrice := showtime.Price * float64(len(seats))

	bookingCode, err := generateBookingCode()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.Add(10 * time.Minute)

	booking := &models.Booking{
		BookingCode: bookingCode,
		UserID:      userID,
		ShowtimeID:  req.ShowtimeID,
		TotalPrice:  totalPrice,
		Status:      "pending",
		ExpiresAt:   &expiresAt,
	}

	bookingSeats := make([]models.BookingSeat, 0, len(seats))

	for _, seat := range seats {
		bookingSeats = append(bookingSeats, models.BookingSeat{
			BookingID:  booking.ID,
			ShowtimeID: req.ShowtimeID,
			SeatID:     seat.ID,
			Price:      showtime.Price,
		})
	}

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txBookingRepo := repository.NewBookingRepository(tx)

		if err := txBookingRepo.CreateBooking(ctx, tx, booking); err != nil {
			return err
		}

		for i := range bookingSeats {
			bookingSeats[i].BookingID = booking.ID
		}

		if err := txBookingRepo.CreateBookingSeats(ctx, tx, bookingSeats); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if isDuplicateEntryError(err) {
			return nil, ErrSeatAlreadyBooked
		}

		return nil, err
	}

	return s.GetBookingByID(ctx, userID, booking.ID)
}

func (s *BookingServiceImpl) GetBookingByID(ctx context.Context, userID uint, id uint) (*response.BookingResponse, error) {
	booking, err := s.Repo.FindBookingByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrBookingNotFound) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}

	if booking.UserID != userID {
		return nil, ErrBookingNotFound
	}

	return toBookingResponse(booking), nil
}

func (s *BookingServiceImpl) GetMyBookings(ctx context.Context, userID uint) ([]response.BookingResponse, error) {
	bookings, err := s.Repo.FindBookingByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]response.BookingResponse, 0, len(bookings))

	for _, booking := range bookings {
		result = append(result, *toBookingResponse(&booking))
	}

	return result, nil
}

func (s *BookingServiceImpl) CancelBooking(ctx context.Context, userID uint, id uint) error {
	booking, err := s.Repo.FindBookingByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrBookingNotFound) {
			return ErrBookingNotFound
		}
		return err
	}

	if booking.UserID != userID {
		return ErrBookingNotFound
	}

	if booking.Status != "pending" {
		return ErrBookingCannotCancel
	}

	if err := s.Repo.CancelBooking(ctx, id); err != nil {
		if errors.Is(err, repository.ErrBookingNotFound) {
			return ErrBookingNotFound
		}

		return err
	}

	return nil
}

func toBookingResponse(booking *models.Booking) *response.BookingResponse {
	return &response.BookingResponse{
		ID:           booking.ID,
		BookingCode:  booking.BookingCode,
		UserID:       booking.UserID,
		ShowtimeID:   booking.ShowtimeID,
		TotalPrice:   booking.TotalPrice,
		Status:       booking.Status,
		ExpiresAt:    booking.ExpiresAt,
		CreatedAt:    booking.CreatedAt,
		UpdatedAt:    booking.UpdatedAt,
		BookingSeats: []response.BookingSeatResponse{},
	}
}

func generateBookingCode() (string, error) {
	buf := make([]byte, 0)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return fmt.Sprintf("BK-%X%", buf), nil
}

func isDuplicateEntryError(err error) bool {
	return err != nil &&
		(errors.Is(err, gorm.ErrDuplicatedKey) ||
			err.Error() == "Error 1062 (23000): Duplicate Entry")
}

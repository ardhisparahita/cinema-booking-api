package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	redisstore "github.com/ardhisparahita/cinema-booking-api/pkg/redis"
	"github.com/go-sql-driver/mysql"
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
	ErrSeatLocked           = errors.New("one or more seats are currently locked")
)

type BookingServiceImpl struct {
	Repo         repository.BookingRepository
	ShowtimeRepo repository.ShowtimeRepository
	SeatRepo     repository.SeatRepository
	DB           *gorm.DB
	SeatLocker   redisstore.SeatLocker
	SeatLockTTL  time.Duration
}

func NewBookingService(repo repository.BookingRepository, showtimeRepo repository.ShowtimeRepository, seatRepo repository.SeatRepository, db *gorm.DB, seatLocker redisstore.SeatLocker, seatLockTTL time.Duration) BookingService {
	return &BookingServiceImpl{
		Repo:         repo,
		ShowtimeRepo: showtimeRepo,
		SeatRepo:     seatRepo,
		DB:           db,
		SeatLocker:   seatLocker,
		SeatLockTTL:  seatLockTTL,
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
		return nil, fmt.Errorf("failed to generate booking code: %w", err)
	}

	locked, err := s.SeatLocker.LockSeats(ctx, req.ShowtimeID, req.SeatIDs, bookingCode, s.SeatLockTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to lock seats: %w", err)
	}

	if !locked {
		return nil, ErrSeatLocked
	}

	unlock := true

	defer func() {
		if unlock {
			_ = s.SeatLocker.UnlockSeats(
				ctx,
				req.ShowtimeID,
				req.SeatIDs,
				bookingCode,
			)
		}
	}()

	now := time.Now()
	expiresAt := now.Add(s.SeatLockTTL)

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

	unlock = false

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

	seatIDs := make([]uint, 0, len(booking.BookingSeats))
	for _, bookingSeat := range booking.BookingSeats {
		seatIDs = append(seatIDs, bookingSeat.SeatID)
	}

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txBookingRepo := repository.NewBookingRepository(tx)

		if err := txBookingRepo.CancelBooking(ctx, id); err != nil {
			if errors.Is(err, repository.ErrBookingNotPending) {
				return ErrBookingCannotCancel
			}
			return err
		}

		if err := txBookingRepo.DeleteBookingSeats(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := s.SeatLocker.UnlockSeats(ctx, booking.ShowtimeID, seatIDs, booking.BookingCode); err != nil {
		log.Printf("failed to unlock redis seats for booking %s: %v", booking.BookingCode, err)
	}

	return nil
}

func (s *BookingServiceImpl) ExpireBooking(ctx context.Context) error {
	now := time.Now()

	bookings, err := s.Repo.FindExpiredPendingBookings(ctx, now)
	if err != nil {
		return err
	}

	for _, booking := range bookings {
		err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			txBookingRepo := repository.NewBookingRepository(tx)
			if err := txBookingRepo.ExpireBooking(ctx, tx, booking.ID); err != nil {
				return err
			}

			if err := txBookingRepo.DeleteBookingSeats(ctx, tx, booking.ID); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to expire booking %d: %w", booking.ID, err)
		}
	}

	return nil
}

func toBookingResponse(booking *models.Booking) *response.BookingResponse {
	bookingSeats := make([]response.BookingSeatResponse, 0, len(booking.BookingSeats))

	for _, bookingSeat := range booking.BookingSeats {
		bookingSeats = append(bookingSeats, response.BookingSeatResponse{
			ID:         bookingSeat.ID,
			SeatID:     bookingSeat.SeatID,
			RowLabel:   bookingSeat.Seat.RowLabel,
			ColNumber:  bookingSeat.Seat.ColNumber,
			SeatType:   bookingSeat.Seat.SeatType,
			ShowtimeID: bookingSeat.ShowtimeID,
			Price:      bookingSeat.Price,
		})
	}

	return &response.BookingResponse{
		ID:          booking.ID,
		BookingCode: booking.BookingCode,
		UserID:      booking.UserID,
		ShowtimeID:  booking.ShowtimeID,
		TotalPrice:  booking.TotalPrice,
		Status:      booking.Status,
		ExpiresAt:   booking.ExpiresAt,
		CreatedAt:   booking.CreatedAt,
		UpdatedAt:   booking.UpdatedAt,
		Seats:       bookingSeats,
	}
}

func generateBookingCode() (string, error) {
	buf := make([]byte, 8)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return fmt.Sprintf("BK-%X", buf), nil
}

func isDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}

	return errors.Is(err, gorm.ErrDuplicatedKey)
}

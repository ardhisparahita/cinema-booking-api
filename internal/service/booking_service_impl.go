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
	appErrors "github.com/ardhisparahita/cinema-booking-api/internal/errors"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	redisstore "github.com/ardhisparahita/cinema-booking-api/pkg/redis"
	"github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type BookingServiceImpl struct {
	Repo         repository.BookingRepository
	ShowtimeRepo repository.ShowtimeRepository
	SeatRepo     repository.SeatRepository
	DB           *gorm.DB
	SeatLocker   redisstore.SeatLocker
	SeatLockTTL  time.Duration
	Tracer       trace.Tracer
}

func NewBookingService(repo repository.BookingRepository, showtimeRepo repository.ShowtimeRepository, seatRepo repository.SeatRepository, db *gorm.DB, seatLocker redisstore.SeatLocker, seatLockTTL time.Duration) BookingService {
	return &BookingServiceImpl{
		Repo:         repo,
		ShowtimeRepo: showtimeRepo,
		SeatRepo:     seatRepo,
		DB:           db,
		SeatLocker:   seatLocker,
		SeatLockTTL:  seatLockTTL,
		Tracer:       otel.Tracer("cinema-booking-api"),
	}
}

func (s *BookingServiceImpl) CreateBooking(
	ctx context.Context,
	userID uint,
	req request.CreateBookingRequest,
) (*response.BookingResponse, error) {

	ctx, span := s.Tracer.Start(
		ctx,
		"BookingService.CreateBooking",
	)
	defer span.End()

	span.SetAttributes(
		attribute.Int64("user.id", int64(userID)),
		attribute.Int64("showtime.id", int64(req.ShowtimeID)),
		attribute.Int("booking.seat_count", len(req.SeatIDs)),
	)

	recordSpanError := func(err error, message string) error {
		if err != nil {
			span.RecordError(err)
		}

		span.SetStatus(codes.Error, message)

		return err
	}

	if len(req.SeatIDs) == 0 {
		return nil, recordSpanError(
			appErrors.ErrInvalidBookingSeats,
			"invalid booking seats",
		)
	}

	seatIDs := uniqueUint(req.SeatIDs)

	if len(seatIDs) != len(req.SeatIDs) {
		return nil, recordSpanError(
			appErrors.ErrDuplicateSeat,
			"duplicate seat selected",
		)
	}

	showtime, err := s.ShowtimeRepo.FindShowtimeByID(
		ctx,
		req.ShowtimeID,
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrShowtimeNotFound) {
			return nil, recordSpanError(
				appErrors.ErrShowtimeNotFound,
				"showtime not found",
			)
		}

		return nil, recordSpanError(
			err,
			"failed to find showtime",
		)
	}

	now := time.Now()

	if !showtime.EndTime.After(now) {
		return nil, recordSpanError(
			appErrors.ErrShowtimeFinished,
			"showtime has already finished",
		)
	}

	seats, err := s.SeatRepo.FindSeatByIDs(
		ctx,
		seatIDs,
	)
	if err != nil {
		return nil, recordSpanError(
			err,
			"failed to find seats",
		)
	}

	if len(seats) != len(seatIDs) {
		return nil, recordSpanError(
			appErrors.ErrSeatNotFound,
			"one or more seats not found",
		)
	}

	for _, seat := range seats {
		if seat.StudioID != showtime.StudioID {
			return nil, recordSpanError(
				appErrors.ErrSeatWrongStudio,
				"seat does not belong to showtime studio",
			)
		}
	}

	bookedSeatIDs, err := s.Repo.FindBookedSeatIDs(
		ctx,
		req.ShowtimeID,
		seatIDs,
	)
	if err != nil {
		return nil, recordSpanError(
			err,
			"failed to check booked seats",
		)
	}

	if len(bookedSeatIDs) > 0 {
		return nil, recordSpanError(
			appErrors.ErrSeatAlreadyBooked,
			"one or more seats are already booked",
		)
	}

	totalPrice := showtime.Price * float64(len(seats))

	span.SetAttributes(
		attribute.Float64("booking.total_price", totalPrice),
	)

	bookingCode, err := generateBookingCode()
	if err != nil {
		return nil, recordSpanError(
			fmt.Errorf("failed to generate booking code: %w", err),
			"failed to generate booking code",
		)
	}

	locked, err := s.SeatLocker.LockSeats(
		ctx,
		req.ShowtimeID,
		seatIDs,
		bookingCode,
		s.SeatLockTTL,
	)
	if err != nil {
		return nil, recordSpanError(
			fmt.Errorf("failed to lock seats: %w", err),
			"failed to execute seat lock",
		)
	}

	if !locked {
		span.SetAttributes(
			attribute.Bool("booking.seat_lock_acquired", false),
		)

		span.AddEvent("seat_lock_denied")

		span.RecordError(appErrors.ErrSeatLocked)
		span.SetStatus(
			codes.Error,
			"seat lock denied",
		)

		return nil, appErrors.ErrSeatLocked
	}

	span.SetAttributes(
		attribute.Bool("booking.seat_lock_acquired", true),
	)

	span.AddEvent("seat_lock_acquired")

	unlock := true

	defer func() {
		if !unlock {
			return
		}

		if err := s.SeatLocker.UnlockSeats(
			ctx,
			req.ShowtimeID,
			seatIDs,
			bookingCode,
		); err != nil {

			span.RecordError(err)

			span.AddEvent("failed_to_unlock_seats")
		}
	}()

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
		bookingSeats = append(
			bookingSeats,
			models.BookingSeat{
				ShowtimeID: req.ShowtimeID,
				SeatID:     seat.ID,
				Price:      showtime.Price,
			},
		)
	}

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txBookingRepo := repository.NewBookingRepository(tx)

		if err := txBookingRepo.CreateBooking(
			ctx,
			tx,
			booking,
		); err != nil {
			return err
		}

		for i := range bookingSeats {
			bookingSeats[i].BookingID = booking.ID
		}

		if err := txBookingRepo.CreateBookingSeats(
			ctx,
			tx,
			bookingSeats,
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {

		if isDuplicateEntryError(err) {

			span.RecordError(err)
			span.SetStatus(
				codes.Error,
				"seat already booked during transaction",
			)

			return nil, appErrors.ErrSeatAlreadyBooked
		}

		return nil, recordSpanError(
			err,
			"failed to create booking transaction",
		)
	}

	span.SetAttributes(
		attribute.Int64("booking.id", int64(booking.ID)),
	)

	unlock = false

	result, err := s.GetBookingByID(
		ctx,
		userID,
		booking.ID,
	)
	if err != nil {
		return nil, recordSpanError(
			err,
			"failed to get created booking",
		)
	}

	return result, nil
}

func (s *BookingServiceImpl) GetBookingByID(ctx context.Context, userID uint, id uint) (*response.BookingResponse, error) {
	booking, err := s.Repo.FindBookingByID(ctx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrBookingNotFound) {
			return nil, appErrors.ErrBookingNotFound
		}
		return nil, err
	}

	if booking.UserID != userID {
		return nil, appErrors.ErrBookingNotFound
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
	var (
		showtimeID  uint
		bookingCode string
		seatIDs     []uint
	)
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txBookingRepo := repository.NewBookingRepository(tx)
		booking, err := txBookingRepo.FindBookingByIDForUpdate(ctx, id)
		if err != nil {
			if errors.Is(err, appErrors.ErrBookingNotFound) {
				return appErrors.ErrBookingNotFound
			}
			return err
		}

		if booking.UserID != userID {
			return appErrors.ErrBookingNotFound
		}

		if booking.Status != "pending" {
			return appErrors.ErrBookingCannotCancel
		}

		showtimeID = booking.ShowtimeID
		bookingCode = booking.BookingCode

		seatIDs = make([]uint, 0, len(booking.BookingSeats))
		for _, bookingSeat := range booking.BookingSeats {
			seatIDs = append(seatIDs, bookingSeat.SeatID)
		}

		if err := txBookingRepo.CancelBooking(ctx, id); err != nil {
			if errors.Is(err, appErrors.ErrBookingNotPending) {
				return appErrors.ErrBookingCannotCancel
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

	if err := s.SeatLocker.UnlockSeats(
		ctx,
		showtimeID,
		seatIDs,
		bookingCode,
	); err != nil {
		log.Printf(
			"failed to unlock redis seats for booking %s: %v",
			bookingCode,
			err,
		)
	}

	return nil

}

func (s *BookingServiceImpl) ExpireBooking(ctx context.Context) error {
	now := time.Now()

	bookings, err := s.Repo.FindExpiredPendingBookings(ctx, now)
	if err != nil {
		return err
	}

	var firstError error

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
			if errors.Is(err, appErrors.ErrBookingNotFound) {
				continue
			}
			if firstError == nil {
				firstError = fmt.Errorf("failed to expire booking %d: %w", booking.ID, err)
			}
			continue
		}

		seatIDs := make([]uint, 0, len(booking.BookingSeats))
		for _, bookingSeat := range booking.BookingSeats {
			seatIDs = append(seatIDs, bookingSeat.SeatID)
		}

		if err := s.SeatLocker.UnlockSeats(ctx, booking.ShowtimeID, seatIDs, booking.BookingCode); err != nil {
			log.Printf("failed to unlock redis seats for expired booking %s: %v", booking.BookingCode, err)
		}
	}

	return firstError
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

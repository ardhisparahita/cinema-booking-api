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
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrPaymentBookingNotFound   = errors.New("booking not found")
	ErrPaymentBookingNotPending = errors.New("booking is not Pending")
	ErrPaymentBookingExpired    = errors.New("booking has expired")
	ErrPaymentAlreadyExist      = errors.New("payment already exist")
	ErrPaymentAlreadyPaid       = errors.New("payment already paid")
	ErrInvalidPaymentMethod     = errors.New("invalid payment method")
	ErrPaymentCannotConfirm     = errors.New("payment cannot be confirmed")
)

type PaymentServiceImpl struct {
	Repo        repository.PaymentRepository
	BookingRepo repository.BookingRepository
	DB          *gorm.DB
	SeatLocker  redisstore.SeatLocker
}

func NewPaymentService(repo repository.PaymentRepository, bookingRepo repository.BookingRepository, db *gorm.DB, seatLocker redisstore.SeatLocker) PaymentService {
	return &PaymentServiceImpl{
		Repo:        repo,
		BookingRepo: bookingRepo,
		DB:          db,
		SeatLocker:  seatLocker,
	}
}

func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, userID uint, req request.CreatePaymentRequest) (*response.PaymentResponse, error) {

	if !isValidPaymentMethod(req.Method) {
		return nil, ErrInvalidPaymentMethod
	}

	var payment *models.Payment

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txBookingRepo := repository.NewBookingRepository(tx)
		txPaymentRepo := repository.NewPaymentRepository(tx)

		booking, err := txBookingRepo.FindBookingByIDForUpdate(ctx, req.BookingID)
		if err != nil {
			if errors.Is(err, repository.ErrBookingNotFound) {
				return ErrPaymentBookingNotFound
			}
			return err
		}

		if booking.UserID != userID {
			return ErrPaymentBookingNotFound
		}

		if booking.Status == "expired" {
			return ErrPaymentBookingExpired
		}

		if booking.Status != "pending" {
			return ErrPaymentBookingNotPending
		}

		if booking.ExpiresAt == nil || !booking.ExpiresAt.After(time.Now()) {
			return ErrPaymentBookingExpired
		}

		existingPayment, err := s.Repo.FindPaymentByBookingID(ctx, booking.ID)
		if err == nil && existingPayment != nil {
			return ErrPaymentAlreadyExist
		}

		if err != nil && !errors.Is(err, repository.ErrPaymentNotFound) {
			return err
		}

		payment := &models.Payment{
			BookingID: booking.ID,
			Method:    req.Method,
			Amount:    booking.TotalPrice,
			Status:    "pending",
		}

		if err := txPaymentRepo.CreatePayment(ctx, tx, payment); err != nil {
			if isDuplicatePaymentError(err) {
				return ErrPaymentAlreadyExist
			}
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return toPaymentResponse(payment), nil
}

func (s *PaymentServiceImpl) GetPaymentByID(ctx context.Context, userID uint, id uint) (*response.PaymentResponse, error) {
	payment, err := s.Repo.FindPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPaymentNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	if payment.Booking.UserID != userID {
		return nil, ErrPaymentNotFound
	}

	return toPaymentResponse(payment), nil
}

func (s *PaymentServiceImpl) ConfirmPayment(ctx context.Context, userID uint, id uint) (*response.PaymentResponse, error) {
	paymentSnapshot, err := s.Repo.FindPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPaymentNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	var payment *models.Payment
	var bookingCode string
	var showTimeID uint
	var seatIDs []uint

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txBookingRepo := repository.NewBookingRepository(tx)
		txPaymentRepo := repository.NewPaymentRepository(tx)

		booking, err := txBookingRepo.FindBookingByIDForUpdate(ctx, paymentSnapshot.ID)
		if err != nil {
			return ErrPaymentBookingNotFound
		}

		lockedPayment, err := txPaymentRepo.FindPaymentByIDForUpdate(ctx, id)
		if err != nil {
			return ErrPaymentNotFound
		}

		if booking.UserID != userID {
			return ErrPaymentNotFound
		}

		if booking.Status == "paid" {
			return ErrPaymentAlreadyPaid
		}

		if booking.Status == "expired" {
			return ErrPaymentBookingExpired
		}

		if booking.Status != "pending" {
			return ErrPaymentBookingNotPending
		}

		if booking.ExpiresAt == nil || !booking.ExpiresAt.After(time.Now()) {
			return ErrPaymentBookingExpired
		}

		if lockedPayment.Status != "pending" {
			return ErrPaymentCannotConfirm
		}

		now := time.Now()

		payment.Status = "paid"
		payment.ProviderRef = generateProviderReference()
		payment.PaidAt = &now

		if err := txPaymentRepo.UpdatePayment(ctx, tx, lockedPayment); err != nil {
			return err
		}

		if err := txBookingRepo.ConfirmBooking(ctx, booking.ID); err != nil {
			return err
		}

		payment = lockedPayment
		bookingCode = booking.BookingCode

		seatIDs = make([]uint, len(booking.BookingSeats))

		for _, bookingSeat := range booking.BookingSeats {
			seatIDs = append(seatIDs, bookingSeat.SeatID)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := s.SeatLocker.UnlockSeats(ctx, showTimeID, seatIDs, bookingCode); err != nil {
		log.Printf("failed to unlock seats after payment %d: %v", id, err)
	}

	return toPaymentResponse(payment), nil
}

func isValidPaymentMethod(method string) bool {
	switch method {
	case "qris", "e_wallet":
		return true
	default:
		return false
	}
}

func isDuplicatePaymentError(err error) bool {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}

	return errors.Is(err, gorm.ErrDuplicatedKey)
}

func toPaymentResponse(payment *models.Payment) *response.PaymentResponse {
	return &response.PaymentResponse{
		ID:          payment.ID,
		BookingID:   payment.BookingID,
		Method:      payment.Method,
		Amount:      payment.Amount,
		Status:      payment.Status,
		ProviderRef: payment.ProviderRef,
		PaidAt:      payment.PaidAt,
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
	}
}

func generateProviderReference() string {
	buf := make([]byte, 6)

	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("SIM-%d", time.Now().UnixNano())
	}

	return fmt.Sprintf("SIM-%X", buf)
}

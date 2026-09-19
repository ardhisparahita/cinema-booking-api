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

func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, userID uint, req *request.CreatePaymentRequest) (*response.PaymentResponse, error) {
	booking, err := s.BookingRepo.FindBookingByID(ctx, req.BookingID)
	if err != nil {
		if errors.Is(err, repository.ErrBookingNotFound) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}

	if booking.UserID != userID {
		return nil, ErrPaymentBookingNotFound
	}

	if booking.Status != "pending" {
		return nil, ErrPaymentBookingNotPending
	}

	if booking.ExpiresAt == nil || !booking.ExpiresAt.After(time.Now()) {
		return nil, ErrPaymentBookingExpired
	}

	if !isValidPaymentMethod(req.Method) {
		return nil, ErrInvalidPaymentMethod
	}

	existingPayment, err := s.Repo.FindPaymentByBookingID(ctx, req.BookingID)
	if err == nil && existingPayment != nil {
		return nil, ErrPaymentAlreadyExist
	}

	if err != nil && !errors.Is(err, repository.ErrPaymentNotFound) {
		return nil, err
	}

	payment := &models.Payment{
		BookingID: booking.ID,
		Method:    req.Method,
		Amount:    booking.TotalPrice,
		Status:    "pending",
	}

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txPaymentRepo := repository.NewPaymentRepository(tx)

		if err := txPaymentRepo.CreatePayment(ctx, tx, payment); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if isDuplicatePaymentError(err) {
			return nil, ErrPaymentAlreadyExist
		}

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

	if payment.Status == "paid" {
		return nil, ErrPaymentAlreadyPaid
	}

	if payment.Booking.Status != "pending" {
		return nil, ErrPaymentBookingNotPending
	}

	if payment.Booking.ExpiresAt == nil || !payment.Booking.ExpiresAt.After(time.Now()) {
		return nil, ErrPaymentBookingExpired
	}

	now := time.Now()

	payment.Status = "paid"
	payment.ProviderRef = generateProviderReference()
	payment.PaidAt = &now

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txPaymentRepo := repository.NewPaymentRepository(tx)

		if err := txPaymentRepo.UpdatePayment(ctx, tx, payment); err != nil {
			return err
		}

		result := tx.WithContext(ctx).Model(&models.Booking{}).Where("id = ?", payment.BookingID).Where("status = ?", "pending").Update("status", "confirmed")

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrPaymentBookingNotPending
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentAlreadyPaid
		}
		return nil, err
	}

	seatIDs := make([]uint, 0, len(payment.Booking.BookingSeats))

	for _, bookingSeat := range payment.Booking.BookingSeats {
		seatIDs = append(seatIDs, bookingSeat.SeatID)
	}

	if err := s.SeatLocker.UnlockSeats(ctx, payment.Booking.ShowtimeID, seatIDs, payment.Booking.BookingCode); err != nil {
		return nil, fmt.Errorf("payment confirmed but failed to unlock seats: %w", err)
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
	buf := make([]byte, 0)

	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("SIM-%d", time.Now().UnixNano())
	}

	return fmt.Sprintf("SIM-%X", buf)
}

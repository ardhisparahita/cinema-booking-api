package test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/pkg/config"
	redisstore "github.com/ardhisparahita/cinema-booking-api/pkg/redis"
	goredis "github.com/redis/go-redis/v9"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestDependencies struct {
	DB             *gorm.DB
	RedisClient    *goredis.Client
	BookingService service.BookingService
	PaymentService service.PaymentService
	BookingRepo    repository.BookingRepository
	PaymentRepo    repository.PaymentRepository
	SeatLocker     redisstore.SeatLocker
}

func SetupTestDependencies(t *testing.T) *TestDependencies {
	t.Helper()

	config.LoadEnvTest()

	ctx := context.Background()

	dsn := os.Getenv("DB_DSN_TEST")
	if dsn == "" {
		t.Skip("DB_DSN_TEST is not set")
	}

	db, err := gorm.Open(mysqlDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	redisDB := 15

	if value := os.Getenv("REDIS_DB_TEST"); value != "" {
		db, err := strconv.Atoi(value)
		if err != nil {
			t.Fatalf("invalid REDIS_DB_TEST: %v", err)
		}

		redisDB = db
	}

	redisClient, err := redisstore.NewClient(
		os.Getenv("REDIS_ADDR_TEST"),
		os.Getenv("REDIS_PASSWORD_TEST"),
		redisDB,
	)
	if err != nil {
		t.Fatalf("failed to connect test redis: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
		_ = redisClient.Close()
	})

	if err := redisClient.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("failed to flush test redis: %v", err)
	}

	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	showtimeRepo := repository.NewShowtimeRepository(db)
	seatRepo := repository.NewSeatRepository(db)
	seatLocker := redisstore.NewSeatLocker(redisClient)

	bookingService := service.NewBookingService(bookingRepo, showtimeRepo, seatRepo, db, seatLocker, 10*time.Minute)

	paymentService := service.NewPaymentService(paymentRepo, bookingRepo, db, seatLocker)

	return &TestDependencies{
		DB:             db,
		RedisClient:    redisClient,
		BookingService: bookingService,
		PaymentService: paymentService,
		BookingRepo:    bookingRepo,
		PaymentRepo:    paymentRepo,
		SeatLocker:     seatLocker,
	}
}

type BookingTestData struct {
	User     models.User
	Theater  models.Theater
	Studio   models.Studio
	Movie    models.Movie
	Showtime models.Showtime
	Seats    []models.Seat
}

func createBookingTestData(t *testing.T, db *gorm.DB) *BookingTestData {
	t.Helper()

	now := time.Now()

	user := models.User{
		Name:        "integration test user",
		Email:       fmt.Sprintf("test-%d@example.com", now.UnixNano()),
		Password:    "test-password",
		PhoneNumber: "123456789",
		Role:        "user",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	theater := models.Theater{
		Name:    fmt.Sprintf("test theater %d", now.UnixNano()),
		City:    "test city",
		Address: "test address",
	}

	if err := db.Create(&theater).Error; err != nil {
		t.Fatalf("failed to create theater: %v", err)
	}

	studio := models.Studio{
		TheaterID: theater.ID,
		Name:      fmt.Sprintf("studio test %d", now.UnixNano()),
		TotalRows: 5,
		TotalCols: 5,
	}

	if err := db.Create(&studio).Error; err != nil {
		t.Fatalf("failed create studio: %v", err)
	}

	movie := models.Movie{
		Title:       fmt.Sprintf("test movie %d", now.UnixNano()),
		Synopsis:    "test synopsis",
		DurationMin: 120,
		Rating:      "SU",
		ReleaseDate: now,
		IsActive:    true,
	}

	if err := db.Create(&movie).Error; err != nil {
		t.Fatalf("failed to create movie: %v", err)
	}

	seat1 := models.Seat{
		StudioID:  studio.ID,
		RowLabel:  "A",
		ColNumber: 1,
		SeatType:  "regular",
	}
	seat2 := models.Seat{
		StudioID:  studio.ID,
		RowLabel:  "A",
		ColNumber: 2,
		SeatType:  "regular",
	}

	if err := db.Create(&seat1).Error; err != nil {
		t.Fatalf("failed to create seat 1: %v", err)
	}
	if err := db.Create(&seat2).Error; err != nil {
		t.Fatalf("failed to create seat 1: %v", err)
	}

	showtime := models.Showtime{
		MovieID:   movie.ID,
		StudioID:  studio.ID,
		StartTime: now.Add(1 * time.Hour),
		EndTime:   now.Add(3 * time.Hour),
		Price:     35000,
	}

	if err := db.Create(&showtime).Error; err != nil {
		t.Fatalf("failed to create showtime: %v", err)
	}

	data := &BookingTestData{
		User:     user,
		Theater:  theater,
		Studio:   studio,
		Movie:    movie,
		Showtime: showtime,
		Seats: []models.Seat{
			seat1,
			seat2,
		},
	}

	t.Cleanup(func() {
		cleanUpBookingTestData(db, data)
	})

	return data
}

func cleanUpBookingTestData(db *gorm.DB, data *BookingTestData) {
	db.Where("showtime_id = ?", data.Showtime.ID).Delete(&models.Booking{})

	db.Delete(&models.Showtime{}, data.Showtime.ID)

	for _, seat := range data.Seats {
		db.Delete(&models.Seat{}, seat.ID)
	}

	db.Delete(&models.Movie{}, data.Movie.ID)
	db.Delete(&models.Studio{}, data.Studio.ID)
	db.Delete(&models.Theater{}, data.Theater.ID)
	db.Delete(&models.User{}, data.User.ID)
}

func TestBookingPaymentFlow(t *testing.T) {
	deps := SetupTestDependencies(t)

	data := createBookingTestData(t, deps.DB)

	ctx := context.Background()

	seatIDs := []uint{
		data.Seats[0].ID,
		data.Seats[1].ID,
	}

	bookingReq := request.CreateBookingRequest{
		ShowtimeID: data.Showtime.ID,
		SeatIDs:    seatIDs,
	}

	booking, err := deps.BookingService.CreateBooking(ctx, data.User.ID, bookingReq)
	if err != nil {
		t.Fatalf("create booking failed: %v", err)
	}

	if booking.Status != "pending" {
		t.Fatalf("expected booking status pending, got %s", booking.Status)
	}

	expectedTotal := data.Showtime.Price * 2

	if booking.TotalPrice != expectedTotal {
		t.Fatalf("expected total price %.2f, got %.2f", expectedTotal, booking.TotalPrice)
	}

	if booking.BookingCode == "" {
		t.Fatal("booking code should not be empty")
	}

	if booking.ExpiresAt == nil {
		t.Fatal("booking expires_at should not be nil")
	}

	for _, seatID := range seatIDs {
		key := fmt.Sprintf("booking:showtime:%d:seat:%d", data.Showtime.ID, seatID)

		value, err := deps.RedisClient.Get(ctx, key).Result()
		if err != nil {
			t.Fatalf("failed to get redis lock for seat %d: %v", seatID, err)
		}

		if value != booking.BookingCode {
			t.Fatalf("expected redis value: %s, got %s", booking.BookingCode, value)
		}
	}

	paymentReq := request.CreatePaymentRequest{
		BookingID: booking.ID,
		Method:    "qris",
	}

	payment, err := deps.PaymentService.CreatePayment(ctx, data.User.ID, paymentReq)
	if err != nil {
		t.Fatalf("create payment failed: %v", err)
	}

	if payment.Status != "pending" {
		t.Fatalf("expected payment status pending, got %s", payment.Status)
	}

	if payment.BookingID != booking.ID {
		t.Fatalf("expected payment booking_id %d, got %d", booking.ID, payment.BookingID)
	}

	if payment.Amount != expectedTotal {
		t.Fatalf("expected payment amount %.2f, got %.2f", expectedTotal, payment.Amount)
	}

	confirmedPayment, err := deps.PaymentService.ConfirmPayment(ctx, data.User.ID, payment.ID)
	if err != nil {
		t.Fatalf("confirm payment failed: %v", err)
	}

	if confirmedPayment.Status != "paid" {
		t.Fatalf("expected payment status paid, got %s", confirmedPayment.Status)
	}

	if confirmedPayment.PaidAt == nil {
		t.Fatal("paid_at should not be nil")
	}

	if confirmedPayment.ProviderRef == "" {
		t.Fatal("provider_ref should not be empty")
	}

	confirmedBooking, err := deps.BookingService.GetBookingByID(ctx, data.User.ID, booking.ID)
	if err != nil {
		t.Fatalf("get confirmed booking failed: %v", err)
	}

	if confirmedBooking.Status != "confirmed" {
		t.Fatalf("expected booking status confirmed, got %s", confirmedBooking.Status)
	}

	if len(confirmedBooking.Seats) != 2 {
		t.Fatalf("expected 2 booking seats, got %d", len(confirmedBooking.Seats))
	}

	for _, seatID := range seatIDs {
		key := fmt.Sprintf("booking:showtime:%d:seat:%d", data.Showtime.ID, seatID)

		exist, err := deps.RedisClient.Exists(ctx, key).Result()
		if err != nil {
			t.Fatalf("failed to check redis key for seat %d: %v", seatID, err)
		}

		if exist != 0 {
			t.Fatalf("redis lock still exist for seat: %d", seatID)
		}
	}
}

func TestConcurrentBookingSameSeat(t *testing.T) {
	deps := SetupTestDependencies(t)
	data := createBookingTestData(t, deps.DB)
	ctx := context.Background()

	user2 := models.User{
		Name:        "concurrent user",
		Email:       fmt.Sprintf("user2-%d@example.com", time.Now().UnixNano()),
		Password:    "test-password",
		PhoneNumber: "123456789",
		Role:        "user",
	}
	if err := deps.DB.Create(&user2).Error; err != nil {
		t.Fatalf("failed to create second user: %v", err)
	}

	t.Cleanup(func() {
		deps.DB.Delete(&models.User{}, user2.ID)
	})
	seatID := data.Seats[0].ID

	bookingReq := request.CreateBookingRequest{
		ShowtimeID: data.Showtime.ID,
		SeatIDs: []uint{
			seatID,
		},
	}

	type result struct {
		Booking *response.BookingResponse
		Err     error
	}

	results := make(chan result, 2)
	start := make(chan struct{})

	go func() {
		<-start

		booking, err := deps.BookingService.CreateBooking(ctx, data.User.ID, bookingReq)

		results <- result{
			Booking: booking,
			Err:     err,
		}
	}()

	go func() {
		<-start

		booking, err := deps.BookingService.CreateBooking(ctx, user2.ID, bookingReq)

		results <- result{
			Booking: booking,
			Err:     err,
		}
	}()

	close(start)

	result1 := <-results
	result2 := <-results

	resultSlice := []result{
		result1,
		result2,
	}

	successCount := 0
	failedCount := 0

	for _, result := range resultSlice {
		if result.Err == nil {
			successCount++

			if result.Booking == nil {
				t.Fatal("successful booking should not be nil")
			}
			if result.Booking.Status != "pending" {
				t.Errorf("expected pending booking, got %s", result.Booking.Status)
			}
		} else {
			failedCount++

			if !errors.Is(result.Err, service.ErrSeatLocked) && !errors.Is(result.Err, service.ErrSeatAlreadyBooked) {
				t.Errorf("unexpected error: %v", result.Err)
			}
		}
	}

	if successCount != 1 {
		t.Fatalf("expected exactly 1 successful booking, got: %d", successCount)
	}
	if failedCount != 1 {
		t.Fatalf("expected exactly 1 failed booking, got: %d", failedCount)
	}

	var count int64

	err := deps.DB.Model(&models.BookingSeat{}).Where("showtime_id = ? AND seat_id = ?", data.Showtime.ID, seatID).Count(&count).Error
	if err != nil {
		t.Fatalf("failed to count booking seats: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected exactly1 booking seat, got: %d", count)
	}
}

func TestCancelBookingFlow(t *testing.T) {
	deps := SetupTestDependencies(t)
	data := createBookingTestData(t, deps.DB)

	ctx := context.Background()

	seatID := data.Seats[0].ID

	req := request.CreateBookingRequest{
		ShowtimeID: data.Showtime.ID,
		SeatIDs:    []uint{seatID},
	}

	booking, err := deps.BookingService.CreateBooking(ctx, data.User.ID, req)
	if err != nil {
		t.Fatalf("create booking failed: %v", err)
	}

	if booking.Status != "pending" {
		t.Fatalf("expected pending, got: %s", booking.Status)
	}

	key := fmt.Sprintf("booking:showtime:%d:seat:%d", data.Showtime.ID, seatID)
	exists, err := deps.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		t.Fatalf("failed checking redis lock :%v", err)
	}

	if exists != 1 {
		t.Fatal("expected redis lock to exist")
	}

	err = deps.BookingService.CancelBooking(ctx, data.User.ID, booking.ID)
	if err != nil {
		t.Fatalf("cancel booking failed: %v", err)
	}

	cancelledBooking, err := deps.BookingService.GetBookingByID(ctx, data.User.ID, booking.ID)
	if err != nil {
		t.Fatalf("get cancelled booking failed: %v", err)
	}

	if cancelledBooking.Status != "cancelled" {
		t.Fatalf("expected cancelled, got %s", cancelledBooking.Status)
	}

	var seatCount int64

	err = deps.DB.Model(&models.BookingSeat{}).Where("booking_id = ?", booking.ID).Count(&seatCount).Error
	if err != nil {
		t.Fatalf("failed counting booking seats: %v", err)
	}

	if seatCount != 0 {
		t.Fatalf("expected 0 booking seats after cancellation, got %d", seatCount)
	}

	exists, err = deps.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		t.Fatalf("failed checking redis after cancellation: %v", err)
	}

	if exists != 0 {
		t.Fatal("redis lock still exists after cancellation")
	}

	booking2, err := deps.BookingService.CreateBooking(ctx, data.User.ID, req)
	if err != nil {
		t.Fatalf("create booking after cancellation failed: %v", err)
	}

	if booking2.Status != "pending" {
		t.Fatalf("expected pending booking, got %s", booking2.Status)
	}
}

func TestExpiredBookingFlow(t *testing.T) {
	deps := SetupTestDependencies(t)
	data := createBookingTestData(t, deps.DB)
	ctx := context.Background()

	shortTTL := 2 * time.Second

	bookingService := service.NewBookingService(
		deps.BookingRepo,
		repository.NewShowtimeRepository(deps.DB),
		repository.NewSeatRepository(deps.DB),
		deps.DB,
		deps.SeatLocker,
		shortTTL,
	)

	seatID := data.Seats[0].ID

	req := request.CreateBookingRequest{
		ShowtimeID: data.Showtime.ID,
		SeatIDs:    []uint{seatID},
	}

	booking, err := bookingService.CreateBooking(ctx, data.User.ID, req)
	if err != nil {
		t.Fatalf("create booking failed: %v", err)
	}

	if booking.Status != "pending" {
		t.Fatalf("expected pending, got %s", booking.Status)
	}

	if booking.ExpiresAt == nil {
		t.Fatal("expires_at should not be nil")
	}

	key := fmt.Sprintf("booking:showtime:%d:seat:%d", data.Showtime.ID, seatID)
	exists, err := deps.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		t.Fatalf("failed checking redis lock: %v", err)
	}

	if exists != 1 {
		t.Fatal("expected redis lock to exist")
	}

	time.Sleep(shortTTL + 1*time.Second)

	exists, err = deps.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		t.Fatalf("failed checking redis after ttl: %v", err)
	}

	if exists != 0 {
		t.Fatal("redis lock should have expired")
	}

	err = bookingService.ExpireBooking(ctx)
	if err != nil {
		t.Fatalf("expire booking failed: %v", err)
	}

	expiredBooking, err := bookingService.GetBookingByID(ctx, data.User.ID, booking.ID)
	if err != nil {
		t.Fatalf("get expired booking failed: %v", err)
	}

	if expiredBooking.Status != "expired" {
		t.Fatalf("expected expired, got: %s", expiredBooking.Status)
	}

	var seatCount int64
	err = deps.DB.Model(&models.BookingSeat{}).Where("booking_id = ?", booking.ID).Count(&seatCount).Error

	if err != nil {
		t.Fatalf("failed counting booking seats: %v", err)
	}

	if seatCount != 0 {
		t.Fatalf("expected 0 booking seats after expiry, got: %d", seatCount)
	}

	booking2, err := bookingService.CreateBooking(ctx, data.User.ID, req)
	if err != nil {
		t.Fatalf("create booking after expiry failed: %v", err)
	}

	if booking2.Status != "pending" {
		t.Fatalf("expected pending booking, got %s", booking2.Status)
	}

}

func TestConcurrentCreatePayment(t *testing.T) {
	deps := SetupTestDependencies(t)
	data := createBookingTestData(t, deps.DB)
	ctx := context.Background()

	bookingReq := request.CreateBookingRequest{
		ShowtimeID: data.Showtime.ID,
		SeatIDs:    []uint{data.Seats[0].ID},
	}

	booking, err := deps.BookingService.CreateBooking(ctx, data.User.ID, bookingReq)
	if err != nil {
		t.Fatalf("failed to create booking: %v", err)
	}

	paymentReq := request.CreatePaymentRequest{
		BookingID: booking.ID,
		Method:    "qris",
	}

	type result struct {
		Payment *response.PaymentResponse
		Err     error
	}

	results := make(chan result, 2)
	start := make(chan struct{})

	go func() {
		<-start

		payment, err := deps.PaymentService.CreatePayment(ctx, data.User.ID, paymentReq)
		results <- result{
			Payment: payment,
			Err:     err,
		}
	}()

	go func() {
		<-start

		payment, err := deps.PaymentService.CreatePayment(ctx, data.User.ID, paymentReq)
		results <- result{
			Payment: payment,
			Err:     err,
		}
	}()

	close(start)

	result1 := <-results
	result2 := <-results

	resultSlice := []result{
		result1, result2,
	}

	successCount := 0
	failedCount := 0

	for _, result := range resultSlice {
		if result.Err == nil {
			successCount++

			if result.Payment == nil {
				t.Fatal("successful payment shouldnot be nil")
			}

			if result.Payment.Status != "pending" {
				t.Errorf("expected payment status pending, got: %s", result.Payment.Status)
			}
		} else {
			failedCount++

			if !errors.Is(result.Err, service.ErrPaymentAlreadyExist) {
				t.Errorf("unexpected error: %v", result.Err)
			}
		}
	}

	if successCount != 1 {
		t.Fatalf("expected exactly 1 successful payment, got %d", successCount)
	}

	if failedCount != 1 {
		t.Fatalf("expected exactly 1 failed payment, got %d", failedCount)
	}

	var paymentCount int64

	err = deps.DB.Model(&models.Payment{}).Where("booking_id = ?", booking.ID).Count(&paymentCount).Error

	if err != nil {
		t.Fatalf("failed to count payments: %v", err)
	}

	if paymentCount != 1 {
		t.Fatalf("expected exactly 1 payment, got: %d", paymentCount)
	}

}

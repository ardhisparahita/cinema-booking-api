package test

import (
	"context"
	"os"
	"testing"

	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	redisstore "github.com/ardhisparahita/cinema-booking-api/pkg/redis"
	goredis "github.com/redis/go-redis/v9"
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

	ctx := context.Background()

	dsn := os.Getenv("")
}
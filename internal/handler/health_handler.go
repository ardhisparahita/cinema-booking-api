package handler

import (
	"github.com/gofiber/fiber/v3"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	DB    *gorm.DB
	Redis *goredis.Client
}

func NewHealthHandler(db *gorm.DB, redis *goredis.Client) *HealthHandler {
	return &HealthHandler{
		DB:    db,
		Redis: redis,
	}
}

func (h *HealthHandler) Liveness(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"service": "cinema-booking-api",
	})
}

func (h *HealthHandler) Readiness(c fiber.Ctx) error {
	databaseStatus := "ok"
	redisStatus := "ok"
	status := "ready"

	sqlDB, err := h.DB.DB()
	if err != nil {
		databaseStatus = "error"
		status = "not ready"

		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":   status,
			"database": databaseStatus,
			"redis":    redisStatus,
		})
	}

	if err := sqlDB.PingContext(c.Context()); err != nil {
		databaseStatus = "error"
		status = "not_ready"
	}

	if err := h.Redis.Ping(c.Context()).Err(); err != nil {
		redisStatus = "error"
		status = "not_ready"
	}

	if status == "not_ready" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":   status,
			"database": databaseStatus,
			"redis":    redisStatus,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   status,
		"database": databaseStatus,
		"redis":    redisStatus,
	})
}

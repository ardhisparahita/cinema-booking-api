package middleware

import (
	"strings"

	"github.com/ardhisparahita/cinema-booking-api/pkg/jwt"
	"github.com/gofiber/fiber/v3"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

func Auth(jwtManager *jwt.Manager) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("authorization")

		if authHeader == "" {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				"authorization header required",
			)
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				"invalid authorization header",
			)
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				"status token required",
			)
		}

		claims, err := jwtManager.ParseAccessToken(token)
		if err != nil {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				"invalid or expired access token",
			)
		}

		c.Locals(UserIDKey, claims.UserID)
		c.Locals(RoleKey, claims.Role)

		return c.Next()
	}
}

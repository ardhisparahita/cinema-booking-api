package routes

import (
	"github.com/ardhisparahita/cinema-booking-api/internal/handler"
	"github.com/ardhisparahita/cinema-booking-api/middleware"
	"github.com/ardhisparahita/cinema-booking-api/pkg/jwt"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, jwtManager *jwt.Manager, genreHandler *handler.GenreHandler) {
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	users := api.Group("/users", middleware.Auth(jwtManager))
	users.Get("/me", userHandler.GetProfile)
	users.Put("/me", userHandler.UpdateProfile)
	users.Delete("/me", userHandler.DeleteProfile)

	genres := api.Group("/genres", middleware.Auth(jwtManager))
	genres.Get("/", genreHandler.GetAll)
	genres.Get("/:id", genreHandler.GetByID)

	genres.Post("/", middleware.Role("admin"), genreHandler.Create)
	genres.Put("/:id", middleware.Role("admin"), genreHandler.Update)
	genres.Delete("/:id", middleware.Role("admin"), genreHandler.Delete)
}

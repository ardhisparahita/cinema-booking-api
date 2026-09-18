package routes

import (
	"github.com/ardhisparahita/cinema-booking-api/internal/handler"
	"github.com/ardhisparahita/cinema-booking-api/middleware"
	"github.com/ardhisparahita/cinema-booking-api/pkg/jwt"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, jwtManager *jwt.Manager, genreHandler *handler.GenreHandler, movieHandler *handler.MovieHandler, theaterHandler *handler.TheaterHandler, studioHandler *handler.StudioHandler, seatHandler *handler.SeatHandler, showtimeHandler *handler.ShowtimeHandler, bookingHandler *handler.BookingHandler) {
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

	movies := api.Group("/movies", middleware.Auth(jwtManager))
	movies.Get("/", movieHandler.GetAll)
	movies.Get("/:id", movieHandler.GetByID)

	movies.Post("/", middleware.Role("admin"), movieHandler.Create)
	movies.Put("/:id", middleware.Role("admin"), movieHandler.Update)
	movies.Delete("/:id", middleware.Role("admin"), movieHandler.Delete)

	theaters := api.Group("/theaters", middleware.Auth(jwtManager))
	theaters.Get("/", theaterHandler.GetAll)
	theaters.Get("/:id", theaterHandler.GetByID)

	theaters.Post("/", middleware.Role("admin"), theaterHandler.Create)
	theaters.Put("/:id", middleware.Role("admin"), theaterHandler.Update)
	theaters.Delete("/:id", middleware.Role("admin"), theaterHandler.Delete)

	studios := api.Group("", middleware.Auth(jwtManager))
	studios.Get("/theaters/:theaterId/studios", studioHandler.GetAll)

	studios.Post("/theaters/:theaterId/studios", middleware.Role("admin"), studioHandler.Create)

	studios.Get("/studios/:id", studioHandler.GetByID)

	studios.Put("/studios/:id", middleware.Role("admin"), studioHandler.Update)
	studios.Delete("/studios/:id", middleware.Role("admin"), studioHandler.Delete)

	seats := api.Group("", middleware.Auth(jwtManager))
	seats.Get("/studios/:studioId/seats", seatHandler.GetAll)

	seats.Post("/studios/:studioId/seats", middleware.Role("admin"), seatHandler.Create)

	seats.Get("/seats/:id", seatHandler.GetByID)

	seats.Put("/seats/:id", middleware.Role("admin"), seatHandler.Update)
	seats.Delete("/seats/:id", middleware.Role("admin"), seatHandler.Delete)

	showtimes := api.Group("showtimes", middleware.Auth(jwtManager))
	showtimes.Get("/", showtimeHandler.GetAll)
	showtimes.Get("/:id", showtimeHandler.GetByID)

	showtimes.Post("/", middleware.Role("admin"), showtimeHandler.Create)
	showtimes.Put("/:id", middleware.Role("admin"), showtimeHandler.Update)
	showtimes.Delete("/:id", middleware.Role("admin"), showtimeHandler.Delete)

	bookings := api.Group("bookings", middleware.Auth(jwtManager))
	bookings.Post("/", bookingHandler.Create)
	bookings.Get("/me", bookingHandler.GetMyBookings)
	bookings.Get("/:id", bookingHandler.GetByID)
	bookings.Post("/:id/cancel", bookingHandler.CancelBooking)
}

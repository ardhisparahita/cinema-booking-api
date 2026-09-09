package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/handler"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	"github.com/ardhisparahita/cinema-booking-api/internal/routes"
	"github.com/ardhisparahita/cinema-booking-api/internal/seeders"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/pkg/config"
	"github.com/ardhisparahita/cinema-booking-api/pkg/database"
	"github.com/ardhisparahita/cinema-booking-api/pkg/jwt"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

func main() {
	config.LoadEnv()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := seeders.SeedAdmin(db); err != nil {
		log.Fatalf("failed run seeder admin: %v", err)
	}

	accessMinutes, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MINUTES"))
	if err != nil {
		log.Fatal("invalid ACCESS_TOKEN_MINUTES")
	}

	refreshDays, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS"))
	if err != nil {
		log.Fatal("invalid REFRESH_TOKEN_DAYS")
	}

	jwtManager := jwt.NewManager(
		os.Getenv("JWT_SECRET"),
		os.Getenv("JWT_ISSUER"),
	)

	accessTTL := time.Duration(accessMinutes) * time.Minute
	refreshTTL := time.Duration(refreshDays) * 24 * time.Hour

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	userRepo := repository.NewUserRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	movieRepo := repository.NewMovieRepository(db)

	authService := service.NewAuthService(
		userRepo,
		*jwtManager,
		accessTTL,
		refreshTTL,
	)
	userService := service.NewUserService(userRepo)
	genreService := service.NewGenreService(genreRepo)
	movieService := service.NewMovieService(movieRepo, db)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	genreHandler := handler.NewGenreHandler(genreService)
	movieHandler := handler.NewMovieHandler(movieService)

	routes.SetupRoutes(
		app,
		authHandler,
		userHandler,
		jwtManager,
		genreHandler,
		movieHandler,
	)

	log.Fatal(app.Listen(":3000"))
}

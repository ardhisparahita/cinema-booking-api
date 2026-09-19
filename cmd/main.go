package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/handler"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	"github.com/ardhisparahita/cinema-booking-api/internal/routes"
	"github.com/ardhisparahita/cinema-booking-api/internal/seeders"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/internal/worker"
	"github.com/ardhisparahita/cinema-booking-api/pkg/config"
	"github.com/ardhisparahita/cinema-booking-api/pkg/database"
	"github.com/ardhisparahita/cinema-booking-api/pkg/jwt"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"

	redisstore "github.com/ardhisparahita/cinema-booking-api/pkg/redis"
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

	seatLockMinutes, err := strconv.Atoi(os.Getenv("SEAT_LOCK_MINUTES"))
	if err != nil {
		log.Fatal("invalid SEAT_LOCK_MINUTES")
	}

	jwtManager := jwt.NewManager(
		os.Getenv("JWT_SECRET"),
		os.Getenv("JWT_ISSUER"),
	)

	accessTTL := time.Duration(accessMinutes) * time.Minute
	refreshTTL := time.Duration(refreshDays) * 24 * time.Hour
	seatLockTTL := time.Duration(seatLockMinutes) * time.Minute

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal("invalid REDIS_DB")
	}

	redisClient, err := redisstore.NewClient(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		redisDB,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer redisClient.Close()

	seatLocker := redisstore.NewSeatLocker(redisClient)
	log.Println("Redis connected successfully")

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	userRepo := repository.NewUserRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	theaterRepo := repository.NewTheaterRepository(db)
	studioRepo := repository.NewStudioRepository(db)
	seatRepo := repository.NewSeatRepository(db)
	showtimeRepo := repository.NewShowtimeRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	authService := service.NewAuthService(
		userRepo,
		*jwtManager,
		accessTTL,
		refreshTTL,
	)
	userService := service.NewUserService(userRepo)
	genreService := service.NewGenreService(genreRepo)
	movieService := service.NewMovieService(movieRepo, db)
	theaterService := service.NewTheaterService(theaterRepo)
	studioService := service.NewStudioService(studioRepo, theaterRepo)
	seatService := service.NewSeatService(seatRepo, studioRepo)
	showtimeService := service.NewShowtimeService(showtimeRepo, movieRepo, studioRepo)
	bookingService := service.NewBookingService(bookingRepo, showtimeRepo, seatRepo, db, seatLocker, seatLockTTL)
	bookingExpiryWorker := worker.NewBookingExpiryWorker(bookingService, 1*time.Minute)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	genreHandler := handler.NewGenreHandler(genreService)
	movieHandler := handler.NewMovieHandler(movieService)
	theaterHandler := handler.NewTheaterHandler(theaterService)
	studioHandler := handler.NewStudioHandler(studioService)
	seatHandler := handler.NewSeatHandler(seatService)
	showtimeHandler := handler.NewShowtimeHandler(showtimeService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	routes.SetupRoutes(
		app,
		authHandler,
		userHandler,
		jwtManager,
		genreHandler,
		movieHandler,
		theaterHandler,
		studioHandler,
		seatHandler,
		showtimeHandler,
		bookingHandler,
	)

	go bookingExpiryWorker.Start(context.Background())

	log.Fatal(app.Listen(":3000"))
}

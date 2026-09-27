package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/handler"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	"github.com/ardhisparahita/cinema-booking-api/internal/routes"
	"github.com/ardhisparahita/cinema-booking-api/internal/seeders"
	"github.com/ardhisparahita/cinema-booking-api/internal/service"
	"github.com/ardhisparahita/cinema-booking-api/internal/worker"
	appMiddleware "github.com/ardhisparahita/cinema-booking-api/middleware"
	"github.com/ardhisparahita/cinema-booking-api/pkg/config"
	"github.com/ardhisparahita/cinema-booking-api/pkg/database"
	"github.com/ardhisparahita/cinema-booking-api/pkg/jwt"
	"github.com/ardhisparahita/cinema-booking-api/pkg/logger"
	"github.com/ardhisparahita/cinema-booking-api/pkg/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	recoverer "github.com/gofiber/fiber/v3/middleware/recover"

	redisstore "github.com/ardhisparahita/cinema-booking-api/pkg/redis"
)

func main() {
	appLogger := logger.New()
	slog.SetDefault(appLogger)

	config.LoadEnv()

	appLogger.Info("application starting")

	db, err := database.ConnectDB()
	if err != nil {
		appLogger.Error(
			"failed to connect database",
			"error", err,
		)
		os.Exit(1)
	}

	if err := seeders.SeedAdmin(db); err != nil {
		appLogger.Error(
			"failed to run admin seeder",
			"error", err,
		)
		os.Exit(1)
	}

	accessMinutes, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MINUTES"))
	if err != nil {
		appLogger.Error(
			"invalid ACCESS_TOKEN_MINUTES",
			"error", err,
		)
		os.Exit(1)
	}

	refreshDays, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS"))
	if err != nil {
		appLogger.Error(
			"invalid REFRESH_TOKEN_DAYS",
			"error", err,
		)
		os.Exit(1)
	}

	seatLockMinutes, err := strconv.Atoi(os.Getenv("SEAT_LOCK_MINUTES"))
	if err != nil {
		appLogger.Error(
			"invalid SEAT_LOCK_MINUTES",
			"error", err,
		)
		os.Exit(1)
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
		appLogger.Error(
			"invalid REDIS_DB",
			"error", err,
		)
		os.Exit(1)
	}

	redisClient, err := redisstore.NewClient(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		redisDB,
	)
	if err != nil {
		appLogger.Error(
			"failed to connect redis",
			"error", err,
			"address", os.Getenv("REDIS_ADDR"),
		)
		os.Exit(1)
	}

	defer redisClient.Close()

	seatLocker := redisstore.NewSeatLocker(redisClient)
	appLogger.Info(
		"redis connected",
		"address", os.Getenv("REDIS_ADDR"),
		"db", redisDB,
	)

	healthHandler := handler.NewHealthHandler(db, redisClient)

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	app.Get("/livez", healthHandler.Liveness)
	app.Get("/readyz", healthHandler.Readiness)

	app.Use(recoverer.New())

	app.Use(requestid.New())
	app.Use(appMiddleware.RequestLogger(appLogger))

	userRepo := repository.NewUserRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	theaterRepo := repository.NewTheaterRepository(db)
	studioRepo := repository.NewStudioRepository(db)
	seatRepo := repository.NewSeatRepository(db)
	showtimeRepo := repository.NewShowtimeRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

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
	paymentService := service.NewPaymentService(paymentRepo, bookingRepo, db, seatLocker)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	genreHandler := handler.NewGenreHandler(genreService)
	movieHandler := handler.NewMovieHandler(movieService)
	theaterHandler := handler.NewTheaterHandler(theaterService)
	studioHandler := handler.NewStudioHandler(studioService)
	seatHandler := handler.NewSeatHandler(seatService)
	showtimeHandler := handler.NewShowtimeHandler(showtimeService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

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
		paymentHandler,
	)

	go bookingExpiryWorker.Start(context.Background())

	appLogger.Info(
		"server starting",
		"address", "3000",
	)

	if err := app.Listen(":3000"); err != nil {
		appLogger.Error(
			"server stopped",
			"error", err,
		)
		os.Exit(1)
	}
}

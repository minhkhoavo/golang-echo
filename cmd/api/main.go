package main

import (
	"fmt"
	"log/slog"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"golang-echo/internal/config"
	"golang-echo/internal/handler"
	appMiddleware "golang-echo/internal/middleware"
	"golang-echo/internal/repository"
	"golang-echo/internal/service"
	appConfig "golang-echo/pkg/config"
	"golang-echo/pkg/utils"
)

func main() {
	// Load configuration
	cfg, err := appConfig.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Initialize logger
	logger := utils.InitLogger(cfg.Logging.Level, cfg.Logging.Format)

	// Initialize database
	db, err := config.InitializeDatabase(cfg.GetDSN())
	if err != nil {
		logger.Error("failed to initialize database", slog.Any("error", err))
		panic(err)
	}
	defer db.Close()

	// Initialize translator
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")

	// Create validator
	validator := utils.NewValidator(trans)
	if err := validator.RegisterAllCustomValidators(); err != nil {
		logger.Error("failed to register custom validators", slog.Any("error", err))
		panic(err)
	}

	// Create JWT Manager
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Duration)

	// Initialize rate limiter
	appMiddleware.InitRateLimiter(cfg.RateLimit.RequestsPerMin, logger)

	// Setup repositories & services
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtManager)
	userHandler := handler.NewUserHandler(userService, validator)

	// Setup Echo
	e := echo.New()
	e.Use(middleware.Recover())
	e.Validator = validator
	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	// Add middlewares
	e.Use(appMiddleware.RequestLoggerMiddleware(logger))
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(appMiddleware.RateLimitMiddleware())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Routes
	apiV1 := e.Group("/api/v1")

	// Public routes
	apiV1.POST("/users", userHandler.CreateUser)
	apiV1.POST("/users/login", userHandler.Login)

	// Protected routes
	protected := apiV1.Group("")
	protected.Use(appMiddleware.JWTMiddleware(jwtManager))
	protected.GET("/users", userHandler.FindAllUsers)
	protected.GET("/users/:id", userHandler.FindUserByID)
	protected.GET("/users/by-email", userHandler.FindUserByEmail)

	logger.Info("Starting server", slog.Int("port", cfg.Server.Port))
	e.Start(fmt.Sprintf(":%d", cfg.Server.Port))
}

package main

import (
	"fmt"
	"log/slog"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/time/rate"

	"golang-echo/internal/config"
	"golang-echo/internal/handler"
	appMiddleware "golang-echo/internal/middleware"
	"golang-echo/internal/repository"
	"golang-echo/internal/service"
	appConfig "golang-echo/pkg/config"
	"golang-echo/pkg/response"
	"golang-echo/pkg/utils"
)

func main() {
	// Load configuration
	cfg, err := appConfig.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Initialize slog.Default() globally
	logger := utils.InitLogger(cfg.Logging.Level, cfg.Logging.Format)
	slog.SetDefault(logger)

	slog.Info("Starting application", slog.String("env", cfg.Server.Env))

	// Initialize database
	db, err := config.InitializeDatabase(cfg.GetDSN())
	if err != nil {
		slog.Error("failed to initialize database", slog.Any("error", err))
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
		slog.Error("failed to register custom validators", slog.Any("error", err))
		panic(err)
	}

	// Create JWT Manager
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Duration)

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
	// e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	// 	LogStatus:    true,
	// 	LogMethod:    true,
	// 	LogURI:       true,
	// 	LogRemoteIP:  true,
	// 	LogRequestID: true,
	// 	LogLatency:   true,
	// 	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
	// 		slog.LogAttrs(
	// 			context.Background(),
	// 			slog.LevelInfo,
	// 			"http_request",
	// 			slog.String("method", v.Method),
	// 			slog.String("uri", v.URI),
	// 			slog.Int("status", v.Status),
	// 			slog.String("remote_ip", v.RemoteIP),
	// 			slog.String("latency", v.Latency.String()),
	// 			slog.String("request_id", v.RequestID),
	// 		)
	// 		return nil
	// 	},
	// }))
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	// Rate limiter
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(
			rate.Limit(float64(cfg.RateLimit.RequestsPerMin) / 60),
		),
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return response.TooManyRequests("RATE_LIMIT_EXCEEDED", "Too many requests. Please try again later.", nil)
		},
	}))

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

	slog.Info("Starting server", slog.Int("port", cfg.Server.Port))
	e.Start(fmt.Sprintf(":%d", cfg.Server.Port))
}

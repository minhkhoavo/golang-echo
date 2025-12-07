package main

import (
	"fmt"
	"log"

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
	// Load configuration using Viper
	cfg, err := appConfig.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database with DSN from config
	db, err := config.InitializeDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize translator (DI - no globals)
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")

	// Create validator with translator (DI)
	validator := utils.NewValidator(trans)
	// Register all custom validators
	if err := validator.RegisterAllCustomValidators(); err != nil {
		log.Fatalf("failed to register custom validators: %v", err)
	}

	// Create JWT Manager (DI)
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Duration)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtManager)
	userHandler := handler.NewUserHandler(userService, validator)

	// Setup routes and start server
	e := echo.New()
	e.Use(middleware.Recover())
	e.Validator = validator
	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())

	apiV1 := e.Group("/api/v1")

	// ===== PUBLIC ROUTES (No authentication required) =====
	publicRoutes := apiV1.Group("")
	publicRoutes.POST("/users", userHandler.CreateUser)
	publicRoutes.POST("/users/login", userHandler.Login)

	// ===== PROTECTED ROUTES (Authentication required) =====
	protectedRoutes := apiV1.Group("")
	protectedRoutes.Use(appMiddleware.JWTMiddleware(jwtManager))

	protectedRoutes.GET("/users", userHandler.FindAllUsers)
	protectedRoutes.GET("/users/:id", userHandler.FindUserByID)
	protectedRoutes.GET("/users/by-email", userHandler.FindUserByEmail)

	log.Printf("Starting server on port %d\n", cfg.Server.Port)
	e.Start(fmt.Sprintf(":%d", cfg.Server.Port))
}

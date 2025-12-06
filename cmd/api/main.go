package main

import (
	"golang-echo/internal/config"
	"golang-echo/internal/handler"
	"golang-echo/internal/repository"
	"golang-echo/internal/service"
	"golang-echo/pkg/utils"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://postgres:root@localhost:5432/golang?sslmode=disable"
	db, err := config.InitializeDatabase(dsn)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Setup routes and start server
	e := echo.New()
	e.Validator = utils.NewValidator()

	apiV1 := e.Group("/api/v1")

	userGroup := apiV1.Group("/users")
	userGroup.GET("", userHandler.FindAllUsers)
	userGroup.GET("/:id", userHandler.FindUserByID)
	userGroup.GET("/by-email", userHandler.FindUserByEmail)
	userGroup.POST("", userHandler.CreateUser)

	e.Start(":8080")
}

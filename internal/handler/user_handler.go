package handler

import (
	"golang-echo/internal/model"
	"golang-echo/internal/service"
	"log"

	"github.com/labstack/echo/v4"
)

type IUserHandler interface {
	CreateUser(c echo.Context) error
	FindAllUsers(c echo.Context) error
}

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) IUserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req model.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Bind error: %v", err)
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		log.Printf("Validation error: %v", err)
		return c.JSON(400, map[string]string{"error": "Validation failed"})
	}
	user, err := h.userService.CreateUser(c.Request().Context(), &req)
	if err != nil {
		log.Printf("Create user error: %v", err)
		return c.JSON(500, map[string]string{"error": "Failed to create user", "details": err.Error()})
	}
	return c.JSON(201, user)
}

func (h *UserHandler) FindAllUsers(c echo.Context) error {
	users, err := h.userService.FindAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to retrieve users"})
	}
	return c.JSON(200, users)
}

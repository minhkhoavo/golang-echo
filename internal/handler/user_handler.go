package handler

import (
	"database/sql"
	"golang-echo/internal/model"
	"golang-echo/internal/service"
	"log"

	"github.com/labstack/echo/v4"
)

type IUserHandler interface {
	CreateUser(c echo.Context) error
	FindAllUsers(c echo.Context) error
	FindUserByID(c echo.Context) error
	FindUserByEmail(c echo.Context) error
}

type userHandler struct {
	userService service.IUserService
}

func (h *userHandler) CreateUser(c echo.Context) error {
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

func (h *userHandler) FindAllUsers(c echo.Context) error {
	users, err := h.userService.FindAllUsers(c.Request().Context())
	if err != nil {
		log.Printf("Find all users error: %v", err)
		return c.JSON(500, map[string]string{"error": "Failed to retrieve users"})
	}
	return c.JSON(200, users)
}

func (h *userHandler) FindUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.userService.FindUserByID(c.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(404, map[string]string{"error": "User not found"})
		}
		log.Printf("Find user by ID error: %v", err)
		return c.JSON(500, map[string]string{"error": "Failed to retrieve user"})
	}
	return c.JSON(200, user)
}

func (h *userHandler) FindUserByEmail(c echo.Context) error {
	email := c.QueryParam("email")
	user, err := h.userService.FindUserByEmail(c.Request().Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(404, map[string]string{"error": "User not found"})
		}
		log.Printf("Find user by email error: %v", err)
		return c.JSON(500, map[string]string{"error": "Failed to retrieve user"})
	}
	return c.JSON(200, user)
}

func NewUserHandler(userService service.IUserService) IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

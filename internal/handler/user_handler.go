package handler

import (
	"golang-echo/internal/model"
	"golang-echo/internal/service"
	"golang-echo/pkg/response"

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
		return response.BadRequest("BIND_ERROR", "Invalid request body", err)
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest("VALIDATION_FAILED", "Validation failed", err)
	}
	user, err := h.userService.CreateUser(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(201, user)
}

func (h *userHandler) FindAllUsers(c echo.Context) error {
	users, err := h.userService.FindAllUsers(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(200, users)
}

func (h *userHandler) FindUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.userService.FindUserByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(200, user)
}

func (h *userHandler) FindUserByEmail(c echo.Context) error {
	email := c.QueryParam("email")
	user, err := h.userService.FindUserByEmail(c.Request().Context(), email)
	if err != nil {
		return err
	}
	return c.JSON(200, user)
}

func NewUserHandler(userService service.IUserService) IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

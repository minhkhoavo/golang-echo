package handler

import (
	"golang-echo/internal/model"
	"golang-echo/internal/service"
	"golang-echo/pkg/request"
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
	return response.Created(c, "SUCCESS", "User created successfully", user)
}

func (h *userHandler) FindAllUsers(c echo.Context) error {
	var pagReq request.PaginationReq
	if err := c.Bind(&pagReq); err != nil {
		return response.BadRequest("BIND_ERROR", "Invalid pagination parameters", err)
	}
	offset, limit, page, pageSize := pagReq.GetQueryParams()
	users, total, err := h.userService.FindAllUsers(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}
	totalPages := int(total+int64(pageSize)-1) / pageSize
	pagination := &response.PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: total,
	}
	return response.ListWithPagination(c, "SUCCESS", "Users retrieved successfully", users, pagination)
}

func (h *userHandler) FindUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.userService.FindUserByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return response.Success(c, "SUCCESS", "User retrieved successfully", user)
}

func (h *userHandler) FindUserByEmail(c echo.Context) error {
	email := c.QueryParam("email")
	user, err := h.userService.FindUserByEmail(c.Request().Context(), email)
	if err != nil {
		return err
	}
	return response.Success(c, "SUCCESS", "User retrieved successfully", user)
}

func NewUserHandler(userService service.IUserService) IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

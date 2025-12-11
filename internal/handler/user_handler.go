package handler

import (
	appMiddleware "golang-echo/internal/middleware"
	"golang-echo/internal/model"
	"golang-echo/internal/service"
	"golang-echo/pkg/request"
	"golang-echo/pkg/response"
	"golang-echo/pkg/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type IUserHandler interface {
	CreateUser(c echo.Context) error
	FindAllUsers(c echo.Context) error
	FindUserByID(c echo.Context) error
	FindUserByEmail(c echo.Context) error
	Login(c echo.Context) error
	GetMyInfo(c echo.Context) error
}

type userHandler struct {
	userService service.IUserService
	validator   *utils.CustomValidator
}

func (h *userHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest("BIND_ERROR", "Invalid request body", err)
	}

	if err := c.Validate(&req); err != nil {
		fieldErrors := h.validator.ExtractValidationErrors(err)
		if len(fieldErrors) == 0 {
			c.Logger().Errorf("Validation error (non-field): %v, Type: %T", err, err)
			return response.BadRequest("VALIDATION_FAILED", "Validation failed", err)
		}
		return response.BadRequestWithFields("VALIDATION_FAILED", "Validation failed", fieldErrors)
	}

	loginResp, err := h.userService.Login(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, "SUCCESS", "Login successful", loginResp)
}

func (h *userHandler) CreateUser(c echo.Context) error {
	var req model.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest("BIND_ERROR", "Invalid request body", err)
	}

	if err := c.Validate(&req); err != nil {
		// Extract field-level validation errors using the validator instance
		fieldErrors := h.validator.ExtractValidationErrors(err)
		// If no field errors were extracted, it means validation failed for some other reason
		if len(fieldErrors) == 0 {
			// Log for debugging
			c.Logger().Errorf("Validation error (non-field): %v, Type: %T", err, err)
			return response.BadRequest("VALIDATION_FAILED", "Validation failed", err)
		}
		return response.BadRequestWithFields("VALIDATION_FAILED", "Validation failed", fieldErrors)
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
	userID, err := strconv.Atoi(id)
	if err != nil {
		return response.BadRequest("INVALID_ID", "Invalid user ID format", err)
	}
	user, err := h.userService.FindUserByID(c.Request().Context(), userID)
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

func (h *userHandler) GetMyInfo(c echo.Context) error {
	userID := appMiddleware.GetUserIDFromContext(c)
	if userID == 0 {
		return response.Unauthorized("INVALID_CONTEXT", "User ID not found in context", nil)
	}

	// Convert int to string using strconv
	user, err := h.userService.FindUserByID(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, "SUCCESS", "User info retrieved successfully", user)
}

func NewUserHandler(userService service.IUserService, validator *utils.CustomValidator) IUserHandler {
	return &userHandler{
		userService: userService,
		validator:   validator,
	}
}

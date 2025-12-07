package service

import (
	"context"
	"errors"
	"golang-echo/internal/model"
	"golang-echo/internal/repository"
	"golang-echo/pkg/response"
	"golang-echo/pkg/utils"
)

type IUserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	FindAllUsers(ctx context.Context, limit int, offset int) ([]*model.User, int64, error)
	FindUserByID(ctx context.Context, email string) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
}
type userService struct {
	userRepo   repository.IUserRepository
	jwtManager *utils.JWTManager
}

func (u *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, response.Internal(err)
	}

	// Create user with hashed password
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
	}
	err = u.userRepo.Create(ctx, user)
	if err != nil {
		// Check for duplicate entry using domain error
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, response.Conflict("EMAIL_ALREADY_REGISTERED", "Email is already registered", err)
		}
		return nil, response.Internal(err)
	}
	return user, nil
}

func (u *userService) FindAllUsers(ctx context.Context, limit int, offset int) ([]*model.User, int64, error) {
	return u.userRepo.FindAll(ctx, limit, offset)
}

func (u *userService) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := u.userRepo.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NotFound("USER_NOT_FOUND", "User not found", err)
		}
		return nil, response.Internal(err)
	}
	return user, nil
}

func (u *userService) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NotFound("USER_NOT_FOUND", "User not found", err)
		}
		return nil, response.Internal(err)
	}
	return user, nil
}

func (u *userService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// Find user by email
	user, err := u.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.Unauthorized("INVALID_CREDENTIALS", "Invalid email or password", err)
		}
		return nil, response.Internal(err)
	}

	// Verify password
	if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, response.Unauthorized("INVALID_CREDENTIALS", "Invalid email or password", err)
	}

	// Generate JWT token
	token, err := u.jwtManager.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, response.Internal(err)
	}

	// Return login response without exposing the hashed password
	user.Password = ""

	return &model.LoginResponse{
		AccessToken: token,
		User:        user,
	}, nil
}

func NewUserService(userRepo repository.IUserRepository, jwtManager *utils.JWTManager) IUserService {
	return &userService{userRepo: userRepo, jwtManager: jwtManager}
}

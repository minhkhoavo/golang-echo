package service

import (
	"context"
	"errors"
	"golang-echo/internal/model"
	"golang-echo/internal/repository"
	"golang-echo/pkg/response"
)

type IUserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	FindAllUsers(ctx context.Context, limit int, offset int) ([]*model.User, int64, error)
	FindUserByID(ctx context.Context, email string) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}
type userService struct {
	userRepo repository.IUserRepository
}

func (u *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Validate input
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, response.BadRequest("VALIDATION_FAILED", "All fields are required", nil)
	}

	// Check password length
	if len(req.Password) < 6 {
		return nil, response.BadRequest("PASSWORD_TOO_SHORT", "Password must be at least 6 characters", nil)
	}

	// Create user directly - let database constraint handle duplicate check
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	err := u.userRepo.Create(ctx, user)
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
	if email == "" {
		return nil, response.BadRequest("VALIDATION_FAILED", "Email cannot be empty", nil)
	}
	user, err := u.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NotFound("USER_NOT_FOUND", "User not found", err)
		}
		return nil, response.Internal(err)
	}
	return user, nil
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

package service

import (
	"context"
	"errors"
	"log/slog"
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
	logger     *slog.Logger
}

func NewUserService(userRepo repository.IUserRepository, jwtManager *utils.JWTManager) IUserService {
	return &userService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     slog.Default(),
	}
}

func (u *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to hash password", slog.String("email", req.Email), slog.Any("error", err))
		return nil, response.Internal(err)
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			u.logger.WarnContext(ctx, "email already registered", slog.String("email", req.Email))
			return nil, response.Conflict("EMAIL_ALREADY_REGISTERED", "Email is already registered", err)
		}
		u.logger.ErrorContext(ctx, "failed to create user", slog.String("email", req.Email), slog.Any("error", err))
		return nil, response.Internal(err)
	}

	u.logger.InfoContext(ctx, "user created successfully", slog.String("email", req.Email), slog.Int("user_id", user.ID))
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
	user, err := u.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			u.logger.WarnContext(ctx, "login failed - user not found", slog.String("email", req.Email))
			return nil, response.Unauthorized("INVALID_CREDENTIALS", "Invalid email or password", err)
		}
		u.logger.ErrorContext(ctx, "failed to find user by email", slog.String("email", req.Email), slog.Any("error", err))
		return nil, response.Internal(err)
	}

	if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
		u.logger.WarnContext(ctx, "login failed - invalid password", slog.String("email", req.Email))
		return nil, response.Unauthorized("INVALID_CREDENTIALS", "Invalid email or password", err)
	}

	token, err := u.jwtManager.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to generate jwt token", slog.String("email", req.Email), slog.Any("error", err))
		return nil, response.Internal(err)
	}

	u.logger.InfoContext(ctx, "login successful", slog.String("email", req.Email), slog.Int("user_id", user.ID))
	user.Password = ""

	return &model.LoginResponse{
		AccessToken: token,
		User:        user,
	}, nil
}

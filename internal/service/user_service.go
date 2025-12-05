package service

import (
	"context"
	"golang-echo/internal/model"
	"golang-echo/internal/repository"
)

type IUserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	FindAllUsers(ctx context.Context) ([]*model.User, error)
}
type userService struct {
	userRepo repository.IUserRepository
}

func (userService *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	err := userService.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userService *userService) FindAllUsers(ctx context.Context) ([]*model.User, error) {
	return userService.userRepo.FindAll(ctx)
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

package service

import (
	"context"
	"golang-echo/internal/model"
	"golang-echo/internal/repository"
)

type IUserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	FindAllUsers(ctx context.Context) ([]*model.User, error)
	FindUserByID(ctx context.Context, email string) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}
type userService struct {
	userRepo repository.IUserRepository
}

func (u *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) FindAllUsers(ctx context.Context) ([]*model.User, error) {
	return u.userRepo.FindAll(ctx)
}

func (u *userService) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	return u.userRepo.FindUserByID(ctx, id)
}

func (u *userService) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return u.userRepo.FindUserByEmail(ctx, email)
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

package repository

import (
	"context"
	"golang-echo/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindAll(ctx context.Context) ([]*model.User, error)
	FindUserByID(ctx context.Context, id string) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	query := `SELECT id, name, email, password, created_at, updated_at FROM users ORDER BY id`
	err := r.db.SelectContext(ctx, &users, query)
	return users, err
}

func (r *userRepository) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`
	var user model.User
	err := r.db.GetContext(ctx, &user, query, id)
	return &user, err
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`
	var user model.User
	err := r.db.GetContext(ctx, &user, query, email)
	return &user, err
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{db: db}
}

package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/morning-glorys/drav-backend/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetByUserEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// get user by email
func (r *userRepository) GetByUserEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, auth_provider, role, created_at 
		FROM users 
		WHERE email = $1
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.AuthProvider,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}

// create user
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, auth_provider, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	if user.Role == "" {
		user.Role = "user"
	}
	err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.PasswordHash, user.AuthProvider, user.Role).
		Scan(&user.ID, &user.CreatedAt)

	return err
}

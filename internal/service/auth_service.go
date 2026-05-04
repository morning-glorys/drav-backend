package service

import (
	"context"

	"github.com/morning-glorys/drav-backend/internal/repository"
)

type AuthService interface {
	GoogleLogin(ctx context.Context, googleToken string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// google login
func (s *authService) GoogleLogin(ctx context.Context, googleToken string) (string, error) {
	return "", nil
	// TODO: implementasi login dengan Google OAuth2
}

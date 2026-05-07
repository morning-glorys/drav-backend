package service

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/pkg/utils"
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
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return "", errors.New("Configuration Is Not Set")
	}

	payload, err := idtoken.Validate(ctx, googleToken, clientID)
	if err != nil {
		return "", errors.New("Invalid Google Token")
	}

	email, emailOk := payload.Claims["email"].(string)
	name, nameOk := payload.Claims["name"].(string)

	if !emailOk || !nameOk {
		return "", errors.New("gagal mengekstrak email dan nama dari akun google")
	}

	user, err := s.userRepo.GetByUserEmail(ctx, email)
	if err != nil {
		newUser := &model.User{
			Name:         name,
			Email:        email,
			AuthProvider: "google",
			Role:         "user",
		}
		err = s.userRepo.CreateUser(ctx, newUser)
		if err != nil {
			return "", errors.New("gagal membuat akun baru")
		}
		user = newUser
	}
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", errors.New("gagal membuat sesi login")
	}

	return token, nil

}

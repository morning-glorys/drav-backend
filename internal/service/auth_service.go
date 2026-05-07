package service

import (
	"context"
	"errors"
	"fmt"
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
		return "", errors.New("google client id is not configured")
	}

	// validate google token
	payload, err := idtoken.Validate(ctx, googleToken, clientID)
	if err != nil {
		return "", fmt.Errorf("invalid google token: %w", err)
	}

	email, emailOk := payload.Claims["email"].(string)
	name, nameOk := payload.Claims["name"].(string)

	if !emailOk || !nameOk {
		return "", errors.New("failed to extract google account data")
	}

	user, err := s.userRepo.GetByUserEmail(ctx, email)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {

			newUser := &model.User{
				Name:         name,
				Email:        email,
				AuthProvider: "google",
				Role:         "user",
			}

			err = s.userRepo.CreateUser(ctx, newUser)
			if err != nil {
				return "", fmt.Errorf("failed to create user: %w", err)
			}

			user = newUser

		} else {
			return "", fmt.Errorf("database error while checking user: %w", err)
		}
	}

	if user.ID == 0 {
		return "", errors.New("invalid user id")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

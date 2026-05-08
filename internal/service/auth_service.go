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

var (
	ErrInvalidGoogleToken   = errors.New("invalid google token")
	ErrUnverifiedGoogleMail = errors.New("google email is not verified")
)

type googleTokenValidator func(ctx context.Context, token string, audience string) (*idtoken.Payload, error)

type AuthService interface {
	GoogleLogin(ctx context.Context, googleToken string) (string, error)
}

type authService struct {
	userRepo            repository.UserRepository
	validateGoogleToken googleTokenValidator
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo:            userRepo,
		validateGoogleToken: idtoken.Validate,
	}
}

// google login
func (s *authService) GoogleLogin(ctx context.Context, googleToken string) (string, error) {

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return "", errors.New("google client id is not configured")
	}

	// validate google token
	payload, err := s.validateGoogleToken(ctx, googleToken, clientID)
	if err != nil {
		return "", errors.Join(ErrInvalidGoogleToken, err)
	}

	email, emailOk := payload.Claims["email"].(string)
	name, nameOk := payload.Claims["name"].(string)
	emailVerified, emailVerifiedOk := payload.Claims["email_verified"].(bool)

	if !emailOk || !nameOk {
		return "", errors.New("failed to extract google account data")
	}

	if !emailVerifiedOk || !emailVerified {
		return "", ErrUnverifiedGoogleMail
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

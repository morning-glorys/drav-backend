package service

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
)

type mockUserRepo struct {
	getByEmailFn func(ctx context.Context, email string) (*model.User, error)
	createUserFn func(ctx context.Context, user *model.User) error
}

func (m *mockUserRepo) GetByUserEmail(ctx context.Context, email string) (*model.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user *model.User) error {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, user)
	}
	return nil
}

func TestGoogleLogin_UserNotFound_CreatesNewUser(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	t.Setenv("JWT_SECRET", "super-secret")

	created := false
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return nil, repository.ErrUserNotFound
		},
		createUserFn: func(ctx context.Context, user *model.User) error {
			created = true
			user.ID = 42
			return nil
		},
	}

	svc := &authService{
		userRepo: repo,
		validateGoogleToken: func(ctx context.Context, token string, audience string) (*idtoken.Payload, error) {
			return &idtoken.Payload{
				Claims: map[string]interface{}{
					"email":          "user@example.com",
					"name":           "User Example",
					"email_verified": true,
				},
			}, nil
		},
	}

	token, err := svc.GoogleLogin(context.Background(), "fake-google-token")
	if err != nil {
		t.Fatalf("GoogleLogin returned unexpected error: %v", err)
	}

	if !created {
		t.Fatal("expected CreateUser to be called")
	}

	if token == "" {
		t.Fatal("expected jwt token, got empty string")
	}
}

func TestGoogleLogin_EmailNotVerified_ReturnsError(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")

	repo := &mockUserRepo{}
	svc := &authService{
		userRepo: repo,
		validateGoogleToken: func(ctx context.Context, token string, audience string) (*idtoken.Payload, error) {
			return &idtoken.Payload{
				Claims: map[string]interface{}{
					"email":          "user@example.com",
					"name":           "User Example",
					"email_verified": false,
				},
			}, nil
		},
	}

	_, err := svc.GoogleLogin(context.Background(), "fake-google-token")
	if !errors.Is(err, ErrUnverifiedGoogleMail) {
		t.Fatalf("expected ErrUnverifiedGoogleMail, got: %v", err)
	}
}

func TestGoogleLogin_InvalidGoogleToken_ReturnsTypedError(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")

	repo := &mockUserRepo{}
	svc := &authService{
		userRepo: repo,
		validateGoogleToken: func(ctx context.Context, token string, audience string) (*idtoken.Payload, error) {
			return nil, errors.New("token invalid")
		},
	}

	_, err := svc.GoogleLogin(context.Background(), "bad-google-token")
	if !errors.Is(err, ErrInvalidGoogleToken) {
		t.Fatalf("expected ErrInvalidGoogleToken, got: %v", err)
	}
}

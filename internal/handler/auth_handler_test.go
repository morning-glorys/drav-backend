package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/service"
)

type mockAuthService struct {
	googleLoginFn func(ctx context.Context, googleToken string) (string, error)
}

func (m *mockAuthService) GoogleLogin(ctx context.Context, googleToken string) (string, error) {
	if m.googleLoginFn != nil {
		return m.googleLoginFn(ctx, googleToken)
	}
	return "", nil
}

func TestGoogleLogin_InvalidToken_ReturnsUnauthorizedGenericMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAuthHandler(&mockAuthService{
		googleLoginFn: func(ctx context.Context, googleToken string) (string, error) {
			return "", service.ErrInvalidGoogleToken
		},
	})

	r := gin.New()
	r.POST("/api/auth/google", h.GoogleLogin)

	body := map[string]string{"id_token": "bad-token"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/google", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if got := w.Body.String(); got != "{\"error\":\"authentication failed\"}" {
		t.Fatalf("unexpected body: %s", got)
	}
}

func TestGoogleLogin_InternalError_ReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAuthHandler(&mockAuthService{
		googleLoginFn: func(ctx context.Context, googleToken string) (string, error) {
			return "", errors.New("db exploded with details")
		},
	})

	r := gin.New()
	r.POST("/api/auth/google", h.GoogleLogin)

	body := map[string]string{"id_token": "some-token"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/google", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if got := w.Body.String(); got != "{\"error\":\"internal server error\"}" {
		t.Fatalf("unexpected body: %s", got)
	}
}

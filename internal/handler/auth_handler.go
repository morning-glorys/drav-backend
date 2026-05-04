package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

// google login handler
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	//TODO: implementasi handler untuk login dengan Google OAuth2
}

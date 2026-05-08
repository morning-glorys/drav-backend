package handler

import (
	"errors"
	"net/http"

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
// @Summary Login with Google
// @Description Authenticate user with Google ID token and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param body body GoogleLoginRequest true "Google login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/google [post]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.IDToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id token is required"})
		return
	}

	token, err := h.authService.GoogleLogin(c.Request.Context(), req.IDToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidGoogleToken) || errors.Is(err, service.ErrUnverifiedGoogleMail) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
	})
}

package handler

import (
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
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
		return
	}

	if req.IDToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Token is required"})
		return
	}

	token, err := h.authService.GoogleLogin(c.Request.Context(), req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "google authentication failed",
		})
		return
	}

	if token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
	})
}

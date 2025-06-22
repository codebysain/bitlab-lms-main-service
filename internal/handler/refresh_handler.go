package handler

import (
	"Internship/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RefreshHandler struct {
	authService service.AuthService
}

func NewRefreshHandler(authService service.AuthService) *RefreshHandler {
	return &RefreshHandler{authService: authService}
}

func (h *RefreshHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing refresh token"})
		return
	}

	// pass the context ➜
	accessToken, newRefreshToken, err :=
		h.authService.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

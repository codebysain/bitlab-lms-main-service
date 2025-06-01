package handler

import (
	"Internship/internal/dto"
	"Internship/internal/entities"
	"Internship/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService // Use the interface, not the pointer to struct
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	accessToken, refreshToken, err := h.authService.Authenticate(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	// 1. Extract role from context (set by middleware)
	roleInterface, exists := c.Get("role")
	if !exists || roleInterface != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 2. Bind input
	var req dto.UserRegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 4. Create user entity
	user := &entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
		Role:     req.Role,
	}

	// 5. Save to DB
	if err := h.authService.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

package handler

import (
	"net/http"

	"Veritasbackend/internal/usecase/auth"
	"Veritasbackend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	loginUseCase         *auth.LoginUseCase
	getCurrentUserUseCase *auth.GetCurrentUserUseCase
}

func NewAuthHandler(loginUseCase *auth.LoginUseCase, getCurrentUserUseCase *auth.GetCurrentUserUseCase) *AuthHandler {
	return &AuthHandler{
		loginUseCase:          loginUseCase,
		getCurrentUserUseCase: getCurrentUserUseCase,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.loginUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generar token JWT
	token, err := jwt.GenerateToken(response.User.ID, response.User.Email, response.TenantID, response.User.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"user":     response.User,
		"tenantId": response.TenantID,
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.getCurrentUserUseCase.Execute(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}


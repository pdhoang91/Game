package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/oden/internal/auth"
	"github.com/yourusername/oden/internal/config"
	"github.com/yourusername/oden/internal/model"
)

// RegisterRequest represents the request to register a new user
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the request to login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the response for auth operations
type AuthResponse struct {
	Success bool   `json:"success"`
	UserID  string `json:"user_id,omitempty"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}

// registerHandler handles user registration
func registerHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Check if username or email already exists
	// This would normally be done via DB call
	// For now, let's assume they don't exist

	// Hash the password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Error hashing password",
		})
		return
	}

	// Create a new user
	userID := uuid.New().String()
	user := model.NewUser(userID, req.Username, req.Email, passwordHash)

	// Save the user to the database
	// This would normally be done via DB call
	// For now, let's assume it was saved successfully

	// Create initial player resources
	resources := model.NewPlayerResources(userID, 1000, 100)

	// Save the resources to the database
	// This would normally be done via DB call

	// Generate a JWT token
	token, err := generateToken(userID, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Error generating token",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		UserID:  userID,
		Token:   token,
	})
}

// loginHandler handles user login
func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Check if user exists and password is correct
	// This would normally be done via DB call
	// For now, let's assume the username is "demo" and password is "password"
	if req.Username != "demo" || req.Password != "password" {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Get user by username
	// This would normally be done via DB call
	// For now, let's assume a user ID
	userID := "user_123456"

	// Update last login time
	// This would normally be done via DB call

	// Generate a JWT token
	token, err := generateToken(userID, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Message: "Error generating token",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		UserID:  userID,
		Token:   token,
	})
}

// generateToken generates a JWT token for the user
func generateToken(userID string, c *gin.Context) (string, error) {
	// Get the config from the context
	cfg, exists := c.Get("config")
	if !exists {
		// If not in the context, use a default secret (should not happen in production)
		return auth.GenerateToken(userID, &config.Config{
			Auth: config.AuthConfig{
				JWTSecret:   "default-secret-change-in-production",
				TokenExpiry: 24,
			},
		})
	}

	return auth.GenerateToken(userID, cfg.(*config.Config))
}

// authMiddleware is a middleware to authenticate requests
func authMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "auth_required",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		// Check if the header is a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "invalid_token",
				"message": "Invalid token format",
			})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := auth.ValidateToken(token, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "invalid_token",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store the user ID in the context
		c.Set("userID", claims.UserID)

		// Store the config in the context
		c.Set("config", cfg)

		// Continue
		c.Next()
	}
}
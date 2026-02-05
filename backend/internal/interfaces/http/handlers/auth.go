package handlers

import (
	"errors"
	"net/http"

	"github.com/blytz/live/backend/internal/application/auth"
	"github.com/blytz/live/backend/internal/domain/user"
	appErrors "github.com/blytz/live/backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	service *auth.Service
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest represents token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         UserDTO `json:"user"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    int     `json:"expires_in"`
}

// UserDTO represents user data transfer object
type UserDTO struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	AvatarURL     string `json:"avatar_url"`
	Role          string `json:"role"`
	EmailVerified bool   `json:"email_verified"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, err.Error()))
		return
	}

	resp, err := h.service.Register(c.Request.Context(), &auth.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusCreated, toAuthResponse(resp))
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, err.Error()))
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusOK, toAuthResponse(resp))
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, err.Error()))
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusOK, toAuthResponse(resp))
}

// GetProfile gets current user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		respondError(c, appErrors.ErrUnauthorizedAccess)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		respondError(c, appErrors.ErrUnauthorizedAccess)
		return
	}

	u, err := h.service.GetUser(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusOK, toUserDTO(u))
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		respondError(c, appErrors.ErrUnauthorizedAccess)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		respondError(c, appErrors.ErrUnauthorizedAccess)
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, err.Error()))
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusOK, gin.H{"message": "password changed successfully"})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	respondJSON(c, http.StatusOK, gin.H{"message": "logged out successfully"})
}

// Helper functions
func toAuthResponse(resp *auth.AuthResponse) *AuthResponse {
	return &AuthResponse{
		User:         toUserDTO(resp.User),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}
}

func toUserDTO(u *user.User) UserDTO {
	return UserDTO{
		ID:            u.ID.String(),
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		AvatarURL:     u.AvatarURL,
		Role:          string(u.Role),
		EmailVerified: u.EmailVerified,
	}
}

func respondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

func respondError(c *gin.Context, err error) {
	var appErr *appErrors.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, gin.H{
			"success": false,
			"error":   appErr.Code,
			"message": appErr.Message,
		})
		return
	}
	
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   "INTERNAL_ERROR",
		"message": err.Error(),
	})
}
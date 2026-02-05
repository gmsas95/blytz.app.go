package auth

import (
	"context"
	"time"

	"github.com/blytz/live/backend/internal/domain/user"
	appErrors "github.com/blytz/live/backend/pkg/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service provides authentication use cases
type Service struct {
	userRepo       user.Repository
	sessionRepo    user.SessionRepository
	tokenManager   user.TokenManager
}

// NewService creates a new auth service
func NewService(userRepo user.Repository, sessionRepo user.SessionRepository, tokenManager user.TokenManager) *Service {
	return &Service{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenManager: tokenManager,
	}
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	Role      user.Role
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string
	Password string
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         *user.User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if user exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrInternal, "failed to check user existence")
	}
	if exists {
		return nil, appErrors.New(appErrors.ErrConflict, "user with this email already exists")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrInternal, "failed to hash password")
	}

	// Create user
	role := req.Role
	if role == "" {
		role = user.RoleBuyer
	}

	u := &user.User{
		ID:            uuid.New(),
		Email:         req.Email,
		PasswordHash:  string(hash),
		Role:          role,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Phone:         req.Phone,
		EmailVerified: false,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrInternal, "failed to create user")
	}

	// Generate tokens
	return s.generateTokens(ctx, u)
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Find user
	u, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if appErrors.IsNotFound(err) {
			return nil, appErrors.New(appErrors.ErrUnauthorized, "invalid email or password")
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, appErrors.New(appErrors.ErrUnauthorized, "invalid email or password")
	}

	// Update last login
	now := time.Now()
	u.LastLoginAt = &now
	s.userRepo.Update(ctx, u)

	// Generate tokens
	return s.generateTokens(ctx, u)
}

// RefreshToken generates new tokens from refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	claims, err := s.tokenManager.Validate(refreshToken)
	if err != nil {
		return nil, appErrors.New(appErrors.ErrUnauthorized, "invalid refresh token")
	}

	// Get user
	u, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	return s.generateTokens(ctx, u)
}

// Logout invalidates a session
func (s *Service) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

// GetUser gets user by ID
func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// ChangePassword changes user password
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(currentPassword)); err != nil {
		return appErrors.New(appErrors.ErrUnauthorized, "current password is incorrect")
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return s.userRepo.Update(ctx, u)
}

// generateTokens creates access and refresh tokens
func (s *Service) generateTokens(ctx context.Context, u *user.User) (*AuthResponse, error) {
	accessToken, refreshToken, err := s.tokenManager.Generate(u.ID, u.Email, u.Role)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrInternal, "failed to generate tokens")
	}

	// Create session
	session := &user.Session{
		ID:           uuid.New(),
		UserID:       u.ID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}
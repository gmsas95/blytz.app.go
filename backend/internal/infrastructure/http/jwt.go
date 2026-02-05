package http

import (
	"errors"
	"time"

	"github.com/blytz/live/backend/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTTokenManager implements user.TokenManager
type JWTTokenManager struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   user.Role `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTTokenManager creates a new JWT token manager
func NewJWTTokenManager(secret string, accessExpiry, refreshExpiry time.Duration) *JWTTokenManager {
	return &JWTTokenManager{
		secret:        []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// Generate generates access and refresh tokens
func (m *JWTTokenManager) Generate(userID uuid.UUID, email string, role user.Role) (accessToken, refreshToken string, err error) {
	// Access token
	accessClaims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}
	
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = access.SignedString(m.secret)
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
	}
	
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refresh.SignedString(m.secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Validate validates a token
func (m *JWTTokenManager) Validate(tokenString string) (*user.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return &user.TokenClaims{
			UserID: claims.UserID,
			Email:  claims.Email,
			Role:   claims.Role,
		}, nil
	}

	return nil, errors.New("invalid token claims")
}
package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/blytz.live.remake/backend/internal/cache"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	tokenBlacklistPrefix  = "token:blacklist:"
	refreshTokenKeyPrefix = "refresh_token:"
	blacklistTTL          = 7 * 24 * time.Hour
)

// CustomClaims represents custom JWT claims
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
	cache         *cache.Cache
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, tokenDuration time.Duration, cache *cache.Cache) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
		cache:         cache,
	}
}

// Generate generates a new JWT token pair (access and refresh)
func (manager *JWTManager) Generate(userID uuid.UUID, email, role string) (accessToken, refreshToken string, err error) {
	accessToken, err = manager.generateToken(userID, email, role, manager.tokenDuration)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshDuration := time.Hour * 24 * 7
	refreshToken, err = manager.generateToken(userID, email, role, refreshDuration)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateToken generates a JWT token with specified duration
func (manager *JWTManager) generateToken(userID uuid.UUID, email, role string, duration time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blytz.live",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// Validate validates a JWT token and returns claims
func (manager *JWTManager) Validate(tokenString string) (*CustomClaims, error) {
	if manager.cache != nil {
		tokenHash := manager.hashToken(tokenString)
		blacklistKey := tokenBlacklistPrefix + tokenHash

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		isBlacklisted, err := manager.cache.Exists(ctx, blacklistKey)
		if err == nil && isBlacklisted {
			return nil, errors.New("token has been revoked")
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(manager.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// Refresh generates a new access token from a valid refresh token
func (manager *JWTManager) Refresh(refreshTokenString string) (string, error) {
	claims, err := manager.Validate(refreshTokenString)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	return manager.generateToken(claims.UserID, claims.Email, claims.Role, manager.tokenDuration)
}

// RevokeToken adds a token to the blacklist
func (manager *JWTManager) RevokeToken(tokenString string, expiration time.Duration) error {
	if manager.cache == nil {
		return nil
	}

	tokenHash := manager.hashToken(tokenString)
	blacklistKey := tokenBlacklistPrefix + tokenHash

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if expiration == 0 {
		expiration = blacklistTTL
	}

	return manager.cache.Set(ctx, blacklistKey, "1", expiration)
}

// IsTokenRevoked checks if a token has been revoked
func (manager *JWTManager) IsTokenRevoked(tokenString string) bool {
	if manager.cache == nil {
		return false
	}

	tokenHash := manager.hashToken(tokenString)
	blacklistKey := tokenBlacklistPrefix + tokenHash

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	isBlacklisted, err := manager.cache.Exists(ctx, blacklistKey)
	return err == nil && isBlacklisted
}

// hashToken creates a hash of the token for storage in Redis
func (manager *JWTManager) hashToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(hash[:])
}

// GetTokenExpiration returns the expiration time of a token
func (manager *JWTManager) GetTokenExpiration(tokenString string) (time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(manager.secretKey), nil
	})

	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return time.Time{}, errors.New("invalid token claims")
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, errors.New("token has no expiration")
	}

	return claims.ExpiresAt.Time, nil
}

// StoreRefreshToken stores a refresh token in Redis for validation
func (manager *JWTManager) StoreRefreshToken(userID uuid.UUID, tokenString string) error {
	if manager.cache == nil {
		return nil
	}

	tokenHash := manager.hashToken(tokenString)
	refreshTokenKey := refreshTokenKeyPrefix + userID.String()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return manager.cache.Set(ctx, refreshTokenKey, tokenHash, time.Hour*24*7)
}

// InvalidateRefreshTokens invalidates all refresh tokens for a user
func (manager *JWTManager) InvalidateRefreshTokens(userID uuid.UUID) error {
	if manager.cache == nil {
		return nil
	}

	refreshTokenKey := refreshTokenKeyPrefix + userID.String()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return manager.cache.Delete(ctx, refreshTokenKey)
}

// ValidateRefreshToken validates a refresh token against stored values
func (manager *JWTManager) ValidateRefreshToken(userID uuid.UUID, tokenString string) error {
	if manager.cache == nil {
		return nil
	}

	tokenHash := manager.hashToken(tokenString)
	refreshTokenKey := refreshTokenKeyPrefix + userID.String()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var storedHash string
	err := manager.cache.Get(ctx, refreshTokenKey, &storedHash)
	if err != nil {
		return errors.New("refresh token not found or expired")
	}

	if storedHash != tokenHash {
		return errors.New("refresh token mismatch")
	}

	return nil
}

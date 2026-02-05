package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleBuyer  Role = "buyer"
	RoleSeller Role = "seller"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID            uuid.UUID
	Email         string
	PasswordHash  string
	Role          Role
	FirstName     string
	LastName      string
	AvatarURL     string
	Phone         string
	EmailVerified bool
	LastLoginAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (u *User) CanBid() bool {
	return u.Role == RoleBuyer || u.Role == RoleAdmin
}

func (u *User) CanSell() bool {
	return u.Role == RoleSeller || u.Role == RoleAdmin
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByToken(ctx context.Context, token string) (*Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type PasswordVerifier interface {
	Hash(password string) (string, error)
	Verify(hash, password string) bool
}

type TokenManager interface {
	Generate(userID uuid.UUID, email string, role Role) (accessToken, refreshToken string, err error)
	Validate(token string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserID uuid.UUID
	Email  string
	Role   Role
}

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Token        string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidToken    = errors.New("invalid token")
	ErrSessionExpired  = errors.New("session expired")
)
package postgres

import (
	"context"
	"errors"

	"github.com/blytz/live/backend/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository implements user.Repository
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	model := toUserModel(u)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	u.ID = model.ID
	u.CreatedAt = model.CreatedAt
	u.UpdatedAt = model.UpdatedAt
	return nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	model := toUserModel(u)
	return r.db.WithContext(ctx).Save(model).Error
}

// GetByID gets user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model User
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return toUserDomain(&model), nil
}

// GetByEmail gets user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var model User
	if err := r.db.WithContext(ctx).First(&model, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return toUserDomain(&model), nil
}

// ExistsByEmail checks if email exists
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Helper functions
func toUserModel(u *user.User) *User {
	return &User{
		BaseModel: BaseModel{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Email:         u.Email,
		PasswordHash:  u.PasswordHash,
		Role:          string(u.Role),
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		AvatarURL:     u.AvatarURL,
		Phone:         u.Phone,
		EmailVerified: u.EmailVerified,
		LastLoginAt:   u.LastLoginAt,
	}
}

func toUserDomain(m *User) *user.User {
	return &user.User{
		ID:            m.ID,
		Email:         m.Email,
		PasswordHash:  m.PasswordHash,
		Role:          user.Role(m.Role),
		FirstName:     m.FirstName,
		LastName:      m.LastName,
		AvatarURL:     m.AvatarURL,
		Phone:         m.Phone,
		EmailVerified: m.EmailVerified,
		LastLoginAt:   m.LastLoginAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
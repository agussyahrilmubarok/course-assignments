package repos

import (
	"context"

	"example.com/internal/account/domain"
	"gorm.io/gorm"
)

//go:generate mockery --name=IUserRepository
type IUserRepository interface {
	FindByID(ctx context.Context, userID string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
	DeleteByID(ctx context.Context, userID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewuserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (s *userRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userRepository) Save(ctx context.Context, user *domain.User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

func (s *userRepository) DeleteByID(ctx context.Context, userID string) error {
	return s.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", userID).Error
}

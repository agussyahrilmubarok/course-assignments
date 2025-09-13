package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IUserRepository
type IUserRepository interface {
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteByID(ctx context.Context, id uint) error
}

type userRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewUserRepository(db *gorm.DB, logger zerolog.Logger) IUserRepository {
	return &userRepository{
		DB:     db,
		Logger: logger,
	}
}

func (r *userRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	if err := r.DB.WithContext(ctx).Find(&users).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to fetch users")
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("user_id", id).Msg("user not found")
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		r.Logger.Error().Err(err).Str("email", email).Msg("user not found by email")
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.DB.WithContext(ctx).Save(user).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to save user")
		return nil, err
	}
	return user, nil
}

func (r *userRepository) DeleteByID(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&domain.User{}, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("user_id", id).Msg("failed to delete user")
		return err
	}
	return nil
}

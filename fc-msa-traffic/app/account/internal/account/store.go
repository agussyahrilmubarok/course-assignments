package account

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type IStore interface {
	FindUserByID(ctx context.Context, userID string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	SaveUser(ctx context.Context, user *User) error
	DeleteUserByID(ctx context.Context, userID string) error
	ExistsUserByEmailIgnoreCase(ctx context.Context, email string) bool
}

type store struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewStore(db *gorm.DB, logger zerolog.Logger) IStore {
	return &store{
		db:     db,
		logger: logger,
	}
}

func (s *store) FindUserByID(ctx context.Context, userID string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("failed to find user by id")
		return nil, err
	}

	s.logger.Info().Str("user_id", userID).Msg("fetching user by id")
	return &user, nil
}

func (s *store) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, "LOWER(email) = ?", strings.ToLower(email)).Error; err != nil {
		s.logger.Error().Err(err).Str("user_email", email).Msg("failed to find user by email")
		return nil, err
	}

	s.logger.Info().Str("user_email", email).Msg("fetching user by email")
	return &user, nil
}

func (s *store) SaveUser(ctx context.Context, user *User) error {
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		s.logger.Error().Err(err).Str("user_email", user.Email).Msg("failed to save user")
		return err
	}

	s.logger.Info().Str("user_email", user.Email).Msg("save user successfully")
	return nil
}

func (s *store) DeleteUserByID(ctx context.Context, userID string) error {
	if err := s.db.WithContext(ctx).Delete(&User{}, "id = ?", userID).Error; err != nil {
		s.logger.Error().Err(err).Msg("failed to delete user")
		return err
	}

	s.logger.Info().Str("user_id", userID).Msg("delete user successfully")
	return nil
}

func (s *store) ExistsUserByEmailIgnoreCase(ctx context.Context, email string) bool {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&User{}).
		Where("LOWER(email) = ?", strings.ToLower(email)).
		Count(&count).Error

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to check if email exists")
		return false
	}

	return count > 0
}

package account

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IStore
type IStore interface {
	WithTx(tx *gorm.DB) IStore
	FindUserByID(ctx context.Context, userID string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	SaveUser(ctx context.Context, user *User) error
	DeleteUserByID(ctx context.Context, userID string) error
	ExistsUserByEmailIgnoreCase(ctx context.Context, email string) bool
}

type store struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewStore(db *gorm.DB, log zerolog.Logger) IStore {
	return &store{
		db:  db,
		log: log,
	}
}

func (s *store) WithTx(tx *gorm.DB) IStore {
	return &store{
		db:  tx,
		log: s.log,
	}
}

func (s *store) FindUserByID(ctx context.Context, userID string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		s.log.Error().Err(err).Msg("Failed to find user by ID")
		return nil, err
	}
	return &user, nil
}

func (s *store) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, "LOWER(email) = ?", strings.ToLower(email)).Error; err != nil {
		s.log.Error().Err(err).Msg("Failed to find user by email")
		return nil, err
	}
	return &user, nil
}

func (s *store) SaveUser(ctx context.Context, user *User) error {
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		s.log.Error().Err(err).Msg("Failed to save user")
		return err
	}
	return nil
}

func (s *store) DeleteUserByID(ctx context.Context, userID string) error {
	if err := s.db.WithContext(ctx).Delete(&User{}, "id = ?", userID).Error; err != nil {
		s.log.Error().Err(err).Msg("Failed to delete user")
		return err
	}
	return nil
}

func (s *store) ExistsUserByEmailIgnoreCase(ctx context.Context, email string) bool {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&User{}).
		Where("LOWER(email) = ?", strings.ToLower(email)).
		Count(&count).Error

	if err != nil {
		s.log.Error().Err(err).Msg("Failed to check if email exists")
		return false
	}
	return count > 0
}

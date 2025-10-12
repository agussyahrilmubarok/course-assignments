package account

import (
	"context"

	"gorm.io/gorm"
)

//go:generate mockery --name=IAccountStore
type IAccountStore interface {
	FindByID(ctx context.Context, userID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, userID string) error
}

type accountStore struct {
	db *gorm.DB
}

func NewAccountStore(db *gorm.DB) IAccountStore {
	return &accountStore{db: db}
}

func (s *accountStore) FindByID(ctx context.Context, userID string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *accountStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *accountStore) Save(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

func (s *accountStore) DeleteByID(ctx context.Context, userID string) error {
	return s.db.WithContext(ctx).Delete(&User{}, "id = ?", userID).Error
}

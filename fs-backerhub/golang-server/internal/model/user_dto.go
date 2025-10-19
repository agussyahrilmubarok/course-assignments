package model

import (
	"time"

	"example.com/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserDTO struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Password   string          `json:"password,omitempty"`
	Role       domain.UserRole `json:"role"`
	Occupation string          `json:"occupation,omitempty"`
	ImageName  string          `json:"image_name,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func (u *UserDTO) FromUser(user *domain.User) {
	if user == nil {
		return
	}

	u.ID = user.ID
	u.Name = user.Name
	u.Email = user.Email
	u.Password = user.Password

	if user.Occupation != nil {
		u.Occupation = *user.Occupation
	}

	if user.ImageName != nil {
		u.ImageName = *user.ImageName
	}

	u.Role = domain.UserRole(user.Role)
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}

func (u *UserDTO) HashPassword(plainText string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *UserDTO) ComparePassword(plainText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	return err == nil
}

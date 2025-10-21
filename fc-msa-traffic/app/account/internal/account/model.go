package account

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *SignUpRequest) ToUser() *User {
	userId := uuid.New().String()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)

	return &User{
		ID:       userId,
		Name:     r.Name,
		Email:    r.Email,
		Password: string(hashPassword),
	}
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type ValidateRequest struct {
	Token string `json:"token" validate:"required"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *UserResponse) FromUser(user *User) {
	r.ID = user.ID
	r.Name = user.Name
	r.Email = user.Email
	r.CreatedAt = user.CreatedAt
	r.UpdatedAt = user.UpdatedAt
}

type AccountResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

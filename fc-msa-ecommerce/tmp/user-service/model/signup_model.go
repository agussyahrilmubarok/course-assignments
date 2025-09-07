package model

import (
	"github.com/go-playground/validator/v10"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *SignUpRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type SignUpResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

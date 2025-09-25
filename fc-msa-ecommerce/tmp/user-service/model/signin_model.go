package model

import "github.com/go-playground/validator/v10"

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *SignInRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type SignInResponse struct {
	Token string `json:"string"`
}

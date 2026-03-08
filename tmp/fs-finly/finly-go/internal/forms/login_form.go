package forms

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (form *LoginForm) Validate() error {
	if err := form.validateRequiredFields(); err != nil {
		return err
	}

	if err := form.validateEmail(); err != nil {
		return err
	}

	if err := form.validatePassword(); err != nil {
		return err
	}

	return nil
}

func (form *LoginForm) validateRequiredFields() error {
	if form.Email == "" || form.Password == "" {
		return errors.New("email and password are required")
	}
	return nil
}

func (form *LoginForm) validateEmail() error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, form.Email)
	if err != nil || !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func (form *LoginForm) validatePassword() error {
	if len(form.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

func (form *LoginForm) ComparePassword(storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(form.Password))
	return err == nil
}

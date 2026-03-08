package forms

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type SignupForm struct {
	Email          string `json:"email" form:"email"`
	Password       string `json:"password" form:"password"`
	RepeatPassword string `json:"repeat_password" form:"repeat_password"`
}

func (form *SignupForm) Validate() error {
	if err := form.validateRequiredFields(); err != nil {
		return err
	}

	if err := form.validateEmail(); err != nil {
		return err
	}

	if err := form.validatePassword(); err != nil {
		return err
	}

	if err := form.validateRepeatPassword(); err != nil {
		return err
	}

	return nil
}

func (form *SignupForm) validateRequiredFields() error {
	if form.Email == "" || form.Password == "" || form.RepeatPassword == "" {
		return errors.New("all fields are required")
	}
	return nil
}

func (form *SignupForm) validateEmail() error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, form.Email)
	if err != nil || !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func (form *SignupForm) validatePassword() error {
	if len(form.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasLetter := false
	hasDigit := false

	for _, char := range form.Password {
		switch {
		case 'A' <= char && char <= 'Z', 'a' <= char && char <= 'z':
			hasLetter = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one number")
	}

	return nil
}

func (form *SignupForm) validateRepeatPassword() error {
	if form.Password != form.RepeatPassword {
		return errors.New("repeat password must match the password")
	}
	return nil
}

func (form *SignupForm) HashPassword() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

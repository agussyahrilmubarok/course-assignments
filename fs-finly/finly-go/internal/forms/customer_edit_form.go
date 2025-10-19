package forms

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
)

type CustomerEditForm struct {
	ID      string `json:"id" form:"id"`
	Name    string `json:"name" form:"name"`
	Email   string `json:"email" form:"email"`
	Phone   string `json:"phone" form:"phone"`
	Address string `json:"address" form:"address"`
}

func (f *CustomerEditForm) Validate() error {
	if err := f.validateRequiredFields(); err != nil {
		return err
	}

	if err := f.validateEmail(); err != nil {
		return err
	}

	if err := f.validatePhone(); err != nil {
		return err
	}

	return nil
}

func (f *CustomerEditForm) validateRequiredFields() error {
	if f.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Name is required")
	}
	if f.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email is required")
	}
	if f.Phone == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Phone is required")
	}
	if f.Address == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Address is required")
	}
	return nil
}

func (f *CustomerEditForm) validateEmail() error {
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	if !regexp.MustCompile(emailRegex).MatchString(f.Email) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
	}
	return nil
}

func (f *CustomerEditForm) validatePhone() error {
	phoneRegex := `^\+?[0-9]{10,15}$` // Allows optional "+" at the start, and 10-15 digits
	if !regexp.MustCompile(phoneRegex).MatchString(f.Phone) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid phone number format")
	}
	return nil
}

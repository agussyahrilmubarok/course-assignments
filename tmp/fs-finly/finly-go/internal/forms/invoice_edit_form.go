package forms

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type InvoiceEditForm struct {
	ID       string  `json:"id" form:"id"`
	Amount   float64 `json:"amount" form:"amount"`
	Date     string  `json:"date" form:"date"`
	Status   string  `json:"status" form:"status"`
	Customer string  `json:"customer" form:"customer"`
}

func (f *InvoiceEditForm) Validate() error {
	if err := f.validateRequiredFields(); err != nil {
		return err
	}

	if err := f.validateAmount(); err != nil {
		return err
	}

	if err := f.validateDate(); err != nil {
		return err
	}

	if err := f.validateStatus(); err != nil {
		return err
	}

	return nil
}

func (f *InvoiceEditForm) validateRequiredFields() error {
	if f.Amount == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Amount is required")
	}
	if f.Date == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Date is required")
	}
	if f.Status == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Status is required")
	}
	if f.Customer == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Customer is required")
	}
	return nil
}

func (f *InvoiceEditForm) validateAmount() error {
	if f.Amount <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Amount must be a positive value")
	}
	return nil
}

func (f *InvoiceEditForm) validateDate() error {
	parsedDate, err := time.Parse("2006-01-02", f.Date)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid date format, expected yyyy/mm/dd")
	}

	f.Date = parsedDate.Format("2006-01-02")
	return nil
}

func (f *InvoiceEditForm) validateStatus() error {
	// Define a list of allowed statuses
	validStatuses := []string{"pending", "paid", "overdue"}
	for _, status := range validStatuses {
		if f.Status == status {
			return nil
		}
	}
	return fiber.NewError(fiber.StatusBadRequest, "Invalid status, allowed values are: pending, paid, overdue")
}

package model

type CreateAddressRequest struct {
	Recipient   string `json:"recipient" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Street      string `json:"street" binding:"required"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Country     string `json:"country" binding:"required"`
	IsDefault   bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	Recipient   *string `json:"recipient"`
	PhoneNumber *string `json:"phone_number"`
	Street      *string `json:"street"`
	City        *string `json:"city"`
	State       *string `json:"state"`
	PostalCode  *string `json:"postal_code"`
	Country     *string `json:"country"`
	IsDefault   *bool   `json:"is_default"`
}

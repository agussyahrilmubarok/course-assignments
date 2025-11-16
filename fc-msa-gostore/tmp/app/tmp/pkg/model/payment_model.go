package model

import "time"

type PaymentModel struct {
	ID         string    `json:"id"`
	RefCode    string    `json:"ref_code"`
	UserID     string    `json:"user_id"`
	OrderID    string    `json:"order_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	PaidAt     time.Time `json:"paid_at"`
	InvoiceUrl string    `json:"invoice_url"`
	ExpiredAt  time.Time `json:"expired_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

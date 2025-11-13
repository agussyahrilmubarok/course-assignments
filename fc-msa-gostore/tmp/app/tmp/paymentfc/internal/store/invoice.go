package store

import "time"

type InvoiceStatus string

const (
	InvoicePending InvoiceStatus = "PENDING"
	InvoicePaid    InvoiceStatus = "PAID"
	InvoiceSettled InvoiceStatus = "SETTLED"
	InvoiceExpired InvoiceStatus = "EXPIRED"
	InvoiceUnknown InvoiceStatus = "UNKNOWN_ENUM_VALUE"
)

type Invoice struct {
	ID         string      `json:"id" bson:"_id,omitempty"`
	InvoiceURL string      `json:"invoice_url" bson:"invoice_url"`
	Invoice    interface{} `json:"invoice" bson:"invoice"` // raw JSON
	CreatedAt  time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" bson:"updated_at"`
}

func (i *Invoice) toInvoiceStatus(status string) InvoiceStatus {
	switch status {
	case "PENDING":
		return InvoicePending
	case "PAID":
		return InvoicePaid
	case "SETTLED":
		return InvoiceSettled
	case "EXPIRED":
		return InvoiceExpired
	default:
		return InvoiceUnknown
	}
}

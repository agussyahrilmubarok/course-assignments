package utils

import (
	"fmt"

	"github.com/phpdave11/gofpdf"
)

type InvoicePDF struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"order_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	PaymentUrl string  `json:"payment_url"`
}

func InvoicePDFGenerator(invoice *InvoicePDF, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "arial")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "[FC] Invoice Details")

	pdf.Ln(20)
	pdf.SetFont("Arial", "", 12)

	pdf.Cell(40, 10, fmt.Sprintf("Payment ID: #%d", invoice.ID))
	pdf.Ln(10)

	pdf.Cell(40, 10, fmt.Sprintf("Order ID: #%d", invoice.OrderID))
	pdf.Ln(10)

	pdf.Cell(40, 10, fmt.Sprintf("Total Amount: Rp%.2f", invoice.Amount))
	pdf.Ln(10)

	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", invoice.Status))
	pdf.Ln(10)

	pdf.Cell(40, 10, fmt.Sprintf("Payment Link: %s", invoice.PaymentUrl))

	return pdf.OutputFileAndClose(outputPath)
}

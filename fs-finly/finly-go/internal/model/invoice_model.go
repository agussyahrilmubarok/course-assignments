package model

import (
	"time"

	"github.com/leekchan/accounting"
)

type Invoice struct {
	ID           string  `bson:"_id,omitempty" json:"id"`
	Amount       float64 `bson:"amount" json:"amount"`
	Date         string  `bson:"date" json:"date"`
	Status       string  `bson:"status" json:"status"`
	Owner        string  `bson:"owner" json:"owner"`
	Customer     string  `bson:"customer" json:"customer"`
	CustomerName string  `json:"customerName"`
}

func (i *Invoice) FormatAmount() string {
	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(int(i.Amount))
}

func (i *Invoice) FormatDate() string {
	parsedDate, err := time.Parse("2006-01-02", i.Date)
	if err != nil {
		return i.Date
	}
	return parsedDate.Format("02 January 2006")
}

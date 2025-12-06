package model

import (
	"time"

	"example.com.backend/internal/domain"
)

type TransactionDTO struct {
	ID            string                   `json:"id"`
	Amount        float64                  `json:"amount"`
	Status        domain.TransactionStatus `json:"status"`
	Reference     string                   `json:"reference"`
	Note          string                   `json:"note,omitempty"`
	UserID        string                   `json:"user_id"`
	UserName      *string                  `json:"user_name,omitempty"`
	CampaignID    string                   `json:"campaign_id"`
	CampaignTitle *string                  `json:"campaign_name,omitempty"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
}

func (t *TransactionDTO) FromTransaction(transaction *domain.Transaction) {
	if transaction == nil {
		return
	}

	t.ID = transaction.ID
	t.Amount = transaction.Amount
	t.Status = domain.TransactionStatus(transaction.Status)
	t.Reference = transaction.Reference
	t.Note = transaction.Note
	t.UserID = transaction.UserID
	t.CampaignID = transaction.CampaignID
	t.CreatedAt = transaction.CreatedAt
	t.UpdatedAt = transaction.UpdatedAt

	if transaction.User.Name != "" {
		t.UserName = &transaction.User.Name
	}

	if transaction.Campaign.Title != "" {
		t.CampaignTitle = &transaction.Campaign.Title
	}
}

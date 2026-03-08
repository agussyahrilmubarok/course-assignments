package payloadV1

import (
	"time"

	"example.com.backend/internal/domain"
)

type TransactionRequest struct {
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	CampaignID string  `json:"campaign_id" binding:"required"`
}

type TransactionResponse struct {
	ID                string                   `json:"id"`
	Amount            float64                  `json:"amount"`
	Status            domain.TransactionStatus `json:"status"`
	Reference         string                   `json:"reference"`
	Note              string                   `json:"note,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	UserID            string                   `json:"user_id"`
	UserName          string                   `json:"user_name"`
	UserEmail         string                   `json:"user_email"`
	CampaignID        string                   `json:"campaign_id"`
	CampaignTitle     string                   `json:"campaign_name"`
	CampaignImageName string                   `json:"campaign_image_name"`
}

func (t *TransactionResponse) FromTransaction(transaction *domain.Transaction) {
	if transaction == nil {
		return
	}

	t.ID = transaction.ID
	t.Amount = transaction.Amount
	t.Status = domain.TransactionStatus(transaction.Status)
	t.Reference = transaction.Reference
	t.Note = transaction.Note
	t.CreatedAt = transaction.CreatedAt
	t.UpdatedAt = transaction.UpdatedAt
	t.UserID = transaction.UserID
	t.UserName = transaction.User.Name
	t.UserEmail = transaction.User.Email
	t.CampaignID = transaction.CampaignID
	t.CampaignTitle = transaction.Campaign.Title

	t.CampaignImageName = ""
	if len(transaction.Campaign.CampaignImages) > 0 {
		for _, img := range transaction.Campaign.CampaignImages {
			if img.IsPrimary {
				t.CampaignImageName = img.ImageName
				break
			}
		}

		if t.CampaignImageName == "" {
			t.CampaignImageName = transaction.Campaign.CampaignImages[0].ImageName
		}
	}
}

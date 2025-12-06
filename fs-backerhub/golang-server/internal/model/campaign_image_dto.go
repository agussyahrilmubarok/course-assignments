package model

import (
	"time"

	"example.com.backend/internal/domain"
)

type CampaignImageDTO struct {
	ID         string    `json:"id"`
	ImageName  string    `json:"image_name"`
	IsPrimary  bool      `json:"is_primary"`
	CampaignID string    `json:"campaign_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *CampaignImageDTO) FromCampaignImage(campaignImage *domain.CampaignImage) {
	if campaignImage == nil {
		return
	}

	c.ID = campaignImage.ID
	c.ImageName = campaignImage.ImageName
	c.IsPrimary = campaignImage.IsPrimary
	c.CampaignID = campaignImage.CampaignID
	c.CreatedAt = campaignImage.CreatedAt
	c.UpdatedAt = campaignImage.UpdatedAt
}

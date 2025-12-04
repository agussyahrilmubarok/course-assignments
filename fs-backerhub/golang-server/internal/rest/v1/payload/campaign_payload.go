package payloadV1

import (
	"mime/multipart"
	"time"

	"example.com.backend/internal/domain"
)

type CampaignRequest struct {
	Title            string  `json:"title" binding:"required,min=3,max=100"`
	ShortDescription string  `json:"short_description" binding:"required,min=10,max=255"`
	Description      string  `json:"description" binding:"omitempty"`
	GoalAmount       float64 `json:"goal_amount" binding:"required,gt=0"`
	Perks            string  `json:"perks" binding:"omitempty,min=3,max=255"`
}

type CampaignResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
	GoalAmount       float64   `json:"goal_amount"`
	CurrentAmount    float64   `json:"current_amount"`
	BackerCount      int64     `json:"backer_count"`
	Perks            string    `json:"perks"`
	Slug             string    `json:"slug"`
	ImageName        string    `json:"image_name,omitempty"`
	UserID           string    `json:"user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (c *CampaignResponse) FromCampaign(campaign *domain.Campaign) {
	if campaign == nil {
		return
	}

	c.ID = campaign.ID
	c.Title = campaign.Title
	c.ShortDescription = campaign.ShortDescription
	c.Description = campaign.Description
	c.GoalAmount = campaign.GoalAmount
	c.CurrentAmount = campaign.CurrentAmount
	c.BackerCount = campaign.BackerCount
	c.Perks = campaign.Perks
	c.Slug = campaign.Slug

	c.ImageName = ""
	if len(campaign.CampaignImages) > 0 {
		for _, img := range campaign.CampaignImages {
			if img.IsPrimary {
				c.ImageName = img.ImageName
				break
			}
		}

		if c.ImageName == "" {
			c.ImageName = campaign.CampaignImages[0].ImageName
		}
	}

	c.UserID = campaign.UserID
	c.CreatedAt = campaign.CreatedAt
	c.UpdatedAt = campaign.UpdatedAt
}

type CampaignDetailResponse struct {
	ID               string                  `json:"id"`
	Title            string                  `json:"title"`
	ShortDescription string                  `json:"short_description"`
	Description      string                  `json:"description"`
	GoalAmount       float64                 `json:"goal_amount"`
	CurrentAmount    float64                 `json:"current_amount"`
	BackerCount      int64                   `json:"backer_count"`
	Perks            string                  `json:"perks"`
	Slug             string                  `json:"slug"`
	ImageName        string                  `json:"image_name,omitempty"`
	UserID           string                  `json:"user_id"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	CampaignImages   []CampaignImageResponse `json:"campaign_images"`
}

func (c *CampaignDetailResponse) FromCampaign(campaign *domain.Campaign) {
	if campaign == nil {
		return
	}

	c.ID = campaign.ID
	c.Title = campaign.Title
	c.ShortDescription = campaign.ShortDescription
	c.Description = campaign.Description
	c.GoalAmount = campaign.GoalAmount
	c.CurrentAmount = campaign.CurrentAmount
	c.BackerCount = campaign.BackerCount
	c.Perks = campaign.Perks
	c.Slug = campaign.Slug

	c.ImageName = ""
	if len(campaign.CampaignImages) > 0 {
		for _, img := range campaign.CampaignImages {
			if img.IsPrimary {
				c.ImageName = img.ImageName
				break
			}
		}

		if c.ImageName == "" {
			c.ImageName = campaign.CampaignImages[0].ImageName
		}
	}

	c.UserID = campaign.UserID
	c.CreatedAt = campaign.CreatedAt
	c.UpdatedAt = campaign.UpdatedAt

	if len(campaign.CampaignImages) > 0 {
		// Fill CampaignImage
		foundPrimary := false
		for _, ci := range campaign.CampaignImages {
			if ci.IsPrimary {
				c.ImageName = ci.ImageName
				foundPrimary = true
				break
			}
		}

		// Fill CampaignImage
		if !foundPrimary {
			last := campaign.CampaignImages[len(campaign.CampaignImages)-1]
			c.ImageName = last.ImageName
		}

		// Fill CampaignImages
		var campaignImages []CampaignImageResponse
		for _, ci := range campaign.CampaignImages {
			var campaignImage CampaignImageResponse
			campaignImage.FromCampaignImage(&ci)
			campaignImages = append(campaignImages, campaignImage)
		}
		c.CampaignImages = campaignImages
	}
}

type CampaignImageRequest struct {
	CampaignID    string                `json:"campaign_id"`
	UserID        string                `json:"user_id"`
	CampaignImage *multipart.FileHeader `json:"campaign_image"`
	IsPrimary     bool                  `json:"is_primary"`
}

type CampaignImageResponse struct {
	ID         string    `json:"id"`
	ImageName  string    `json:"image_name"`
	IsPrimary  bool      `json:"is_primary"`
	CampaignID string    `json:"campaign_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *CampaignImageResponse) FromCampaignImage(campaignImage *domain.CampaignImage) {
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

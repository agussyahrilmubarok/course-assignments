package model

import (
	"fmt"
	"time"

	"example.com.backend/internal/domain"
	"github.com/gosimple/slug"
	"github.com/leekchan/accounting"
)

type CampaignDetailDTO struct {
	ID               string             `json:"id"`
	Title            string             `json:"title"`
	ShortDescription string             `json:"short_description"`
	Description      string             `json:"description"`
	GoalAmount       float64            `json:"goal_amount"`
	CurrentAmount    float64            `json:"current_amount"`
	BackerCount      int64              `json:"backer_count"`
	Perks            string             `json:"perks"`
	Slug             string             `json:"slug"`
	ImageName        string             `json:"image_name"`
	UserID           string             `json:"user_id"`
	CampaignImages   []CampaignImageDTO `json:"campaign_images"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

func (c *CampaignDetailDTO) FromCampaign(campaign *domain.Campaign) {
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
		var campaignImages []CampaignImageDTO
		for _, ci := range campaign.CampaignImages {
			var campaignImage CampaignImageDTO
			campaignImage.FromCampaignImage(&ci)
			campaignImages = append(campaignImages, campaignImage)
		}

		c.CampaignImages = campaignImages
	}
}

func (c *CampaignDetailDTO) GenerateSlug() {
	timeUnix := time.Now().Unix()
	slugCandidate := fmt.Sprintf("%s-%s-%d", c.Title, c.UserID, timeUnix)
	c.Slug = slug.Make(slugCandidate)
}

func (c *CampaignDetailDTO) GoalAmountIDR() string {
	return c.formatIDR(int(c.GoalAmount))
}

func (c *CampaignDetailDTO) CurrentAmountIDR() string {
	return c.formatIDR(int(c.CurrentAmount))
}

func (c *CampaignDetailDTO) formatIDR(amount int) string {
	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(amount)
}

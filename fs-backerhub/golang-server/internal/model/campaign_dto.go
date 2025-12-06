package model

import (
	"fmt"
	"time"

	"example.com.backend/internal/domain"
	"github.com/gosimple/slug"
	"github.com/leekchan/accounting"
)

type CampaignDTO struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
	GoalAmount       float64   `json:"goal_amount"`
	CurrentAmount    float64   `json:"current_amount"`
	BackerCount      int64     `json:"backer_count"`
	Perks            string    `json:"perks"`
	Slug             string    `json:"slug"`
	UserID           string    `json:"user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (c *CampaignDTO) FromCampaign(campaign *domain.Campaign) {
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
	c.UserID = campaign.UserID
	c.CreatedAt = campaign.CreatedAt
	c.UpdatedAt = campaign.UpdatedAt
}

func (c *CampaignDTO) GenerateSlug() {
	timeUnix := time.Now().Unix()
	slugCandidate := fmt.Sprintf("%s-%s-%d", c.Title, c.UserID, timeUnix)
	c.Slug = slug.Make(slugCandidate)
}

func (c *CampaignDTO) GoalAmountIDR() string {
	return c.formatIDR(int(c.GoalAmount))
}

func (c *CampaignDTO) CurrentAmountIDR() string {
	return c.formatIDR(int(c.CurrentAmount))
}

func (c *CampaignDTO) formatIDR(amount int) string {
	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(amount)
}

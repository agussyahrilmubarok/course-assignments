package domain

import (
	"time"
)

type CampaignImage struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ImageName  string    `gorm:"type:varchar(255);not null" json:"image_name"`
	IsPrimary  bool      `gorm:"type:boolean;not null;default:false" json:"is_primary"`
	CampaignID string    `gorm:"type:uuid;not null" json:"campaign_id"`
	Campaign   *Campaign `gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"campaign,omitempty"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
}

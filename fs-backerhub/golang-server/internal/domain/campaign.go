package domain

import "time"

type Campaign struct {
	ID               string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title            string    `gorm:"type:varchar(255);not null" json:"title"`
	ShortDescription string    `gorm:"type:varchar(255);not null" json:"short_description"`
	Description      string    `gorm:"type:text" json:"description,omitempty"`
	GoalAmount       float64   `gorm:"type:numeric(10,2);not null" json:"goal_amount"`
	CurrentAmount    float64   `gorm:"type:numeric(10,2);not null" json:"current_amount"`
	BackerCount      int64     `gorm:"type:bigint;not null" json:"backer_count"`
	Perks            string    `gorm:"type:text;not null" json:"perks"`
	Slug             string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	CreatedAt        time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	UserID         string          `gorm:"type:uuid;not null" json:"user_id"`
	User           *User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	CampaignImages []CampaignImage `gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"campaign_images,omitempty"`
}

type CampaignImage struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ImageName string    `gorm:"type:varchar(255);not null" json:"image_name"`
	IsPrimary bool      `gorm:"type:boolean;not null;default:false" json:"is_primary"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	CampaignID string    `gorm:"type:uuid;not null" json:"campaign_id"`
	Campaign   *Campaign `gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"campaign,omitempty"`
}

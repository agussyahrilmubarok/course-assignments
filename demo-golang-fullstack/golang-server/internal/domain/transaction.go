package domain

import (
	"time"
)

type Transaction struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Amount     float64   `gorm:"type:numeric(10,2);not null" json:"amount"`
	Status     string    `gorm:"type:varchar(50);not null" json:"status"`
	Note       string    `gorm:"type:text" json:"note,omitempty"`
	Reference  string    `gorm:"type:varchar(255);not null;uniqueIndex;column:reference" json:"reference"`
	UserID     string    `gorm:"type:uuid;not null" json:"user_id"`
	User       *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	CampaignID string    `gorm:"type:uuid;not null" json:"campaign_id"`
	Campaign   *Campaign `gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"campaign,omitempty"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
}

type TransactionStatus string

const (
	StatusPending  TransactionStatus = "PENDING"
	StatusPaid     TransactionStatus = "PAID"
	StatusFailed   TransactionStatus = "FAILED"
	StatusExpired  TransactionStatus = "EXPIRED"
	StatusCanceled TransactionStatus = "CANCELED"
)

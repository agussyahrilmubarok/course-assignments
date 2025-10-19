package domain

import (
	"time"
)

type User struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Email      string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password   string    `gorm:"type:varchar(255);not null" json:"password"`
	Role       string    `gorm:"type:varchar(50);not null" json:"role"`
	Occupation *string   `gorm:"type:varchar(255)" json:"occupation,omitempty"`
	ImageName  *string   `gorm:"type:varchar(255)" json:"image_name,omitempty"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
}

type UserRole string

const (
	RoleAdmin UserRole = "ROLE_ADMIN"
	RoleUser  UserRole = "ROLE_USER"
)

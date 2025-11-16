package store

import (
	"time"
)

// UserRole defines available roles for a user.
type UserRole string

const (
	RoleAdmin UserRole = "ROLE_ADMIN"
	RoleUser  UserRole = "ROLE_USER"
)

// User represents a user document in MongoDB.
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password,omitempty" bson:"password"`
	Role      UserRole  `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

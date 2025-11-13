package store

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product document in MongoDB.
type Product struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Price     float64   `json:"price" bson:"price"`
	Stock     int       `json:"stock" bson:"stock"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func NewProduct(name string, price float64, stock int) *Product {
	return &Product{
		ID:        uuid.New().String(),
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

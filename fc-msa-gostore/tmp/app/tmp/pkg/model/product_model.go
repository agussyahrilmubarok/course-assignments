package model

import "time"

type ProductModel struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"gte=0"`
	Stock int     `json:"stock" validate:"gte=0"`
}

type UpdateProductRequest struct {
	ID    *string  `json:"id,omitempty"`
	Name  *string  `json:"name,omitempty"`
	Price *float64 `json:"price,omitempty"`
	Stock *int     `json:"stock,omitempty"`
}

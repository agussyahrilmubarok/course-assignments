package domain

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"              json:"id"`
	Name         string    `gorm:"type:varchar(100)"       json:"name"`
	Email        string    `gorm:"uniqueIndex;not null"    json:"email"`
	PasswordHash string    `gorm:"type:varchar(255)"       json:"-"`
	Role         string    `gorm:"type:varchar(20);default:'user'" json:"role"` // user, admin
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Addresses []Address `json:"addresses,omitempty"`
	Orders    []Order   `json:"orders,omitempty"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey"          json:"id"`
	Name        string    `gorm:"type:varchar(100)"   json:"name"`
	Description string    `gorm:"type:text"           json:"description"`
	Price       float64   `gorm:"not null"            json:"price"`
	Stock       int       `gorm:"not null"            json:"stock"`
	ImageURL    string    `gorm:"type:text"           json:"image_url,omitempty"`
	CategoryID  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []Tag    `gorm:"many2many:product_tags" json:"tags,omitempty"`
	Reviews  []Review `json:"reviews,omitempty"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey"        json:"id"`
	Name      string    `gorm:"unique;not null"   json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Products []Product `json:"products,omitempty"`
}

type Tag struct {
	ID        uint      `gorm:"primaryKey"        json:"id"`
	Name      string    `gorm:"unique;not null"   json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Products []Product `gorm:"many2many:product_tags" json:"products,omitempty"`
}

type Cart struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null"   json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Items []CartItem `json:"items,omitempty"`
}

type CartItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CartID    uint    `gorm:"not null"   json:"cart_id"`
	ProductID uint    `gorm:"not null"   json:"product_id"`
	Quantity  int     `gorm:"not null"   json:"quantity"`
	Product   Product `json:"product,omitempty"`
}

type Order struct {
	ID          uint      `gorm:"primaryKey"      json:"id"`
	UserID      uint      `gorm:"not null"        json:"user_id"`
	AddressID   uint      `gorm:"not null"        json:"address_id"`
	TotalAmount float64   `gorm:"not null"        json:"total_amount"`
	Status      string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, paid, shipped, cancelled
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Items   []OrderItem `json:"items,omitempty"`
	Address Address     `json:"address,omitempty"`
	Payment Payment     `json:"payment,omitempty"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null"   json:"order_id"`
	ProductID uint    `gorm:"not null"   json:"product_id"`
	Quantity  int     `gorm:"not null"   json:"quantity"`
	Price     float64 `gorm:"not null"   json:"price"` // Snapshot of product price at purchase

	Product Product `json:"product,omitempty"`
}

type Payment struct {
	ID          uint       `gorm:"primaryKey"     json:"id"`
	OrderID     uint       `gorm:"uniqueIndex"    json:"order_id"`
	Method      string     `gorm:"type:varchar(50)" json:"method"`                   // credit_card, transfer, etc.
	Status      string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, success, failed
	ReferenceID string     `gorm:"type:varchar(100)" json:"reference_id"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Address struct {
	ID          uint   `gorm:"primaryKey"           json:"id"`
	UserID      uint   `gorm:"not null"             json:"user_id"`
	Recipient   string `gorm:"type:varchar(100)"    json:"recipient"`
	PhoneNumber string `gorm:"type:varchar(20)"     json:"phone_number"`
	Street      string `gorm:"type:text"            json:"street"`
	City        string `gorm:"type:varchar(100)"    json:"city"`
	State       string `gorm:"type:varchar(100)"    json:"state"`
	PostalCode  string `gorm:"type:varchar(20)"     json:"postal_code"`
	Country     string `gorm:"type:varchar(100)"    json:"country"`
	IsDefault   bool   `gorm:"default:false"        json:"is_default"`
}

type Review struct {
	ID        uint      `gorm:"primaryKey"     json:"id"`
	UserID    uint      `gorm:"not null"       json:"user_id"`
	ProductID uint      `gorm:"not null"       json:"product_id"`
	Rating    int       `gorm:"not null"       json:"rating"` // 1 to 5
	Comment   string    `gorm:"type:text"      json:"comment"`
	CreatedAt time.Time `json:"created_at"`

	User    User    `json:"user,omitempty"`
	Product Product `json:"product,omitempty"`
}

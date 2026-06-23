package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      string         `json:"user_id" gorm:"not null;type:uuid"`
	Status      OrderStatus    `json:"status" gorm:"default:pending"`
	TotalAmount float64        `json:"total_amount" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// relationship
	User       User        `json:"user"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID   string         `json:"order_id" gorm:"not null;type:uuid"`
	ProductID string         `json:"product_id" gorm:"not null;type:uuid"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	Price     float64        `json:"price" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// relationships
	Order   Order   `json:"-"`
	Product Product `json:"product"`
}

type Cart struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string         `json:"user_id" gorm:"uniqueIndex;not null;type:uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// relationship
	CartItem []CartItem `json:"cart_item"`
}

type CartItem struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CartID    string         `json:"cart_id" gorm:"not null;type:uuid"`
	ProductID string         `json:"product_id" gorm:"not null;type:uuid"`
	Quantity  uint           `json:"quantity" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// relationship
	Cart    Cart    `json:"-"`
	Product Product `json:"product"`
}

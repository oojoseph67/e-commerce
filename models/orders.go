package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	Id          uint           `json:"id" gorm:"primaryKey"`
	UserId      uint           `json:"user_id" gorm:"not null"`
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
	Id        uint           `json:"id" gorm:"primaryKey"`
	OrderId   uint           `json:"order_id" gorm:"not null"`
	ProductId uint           `json:"product_id" gorm:"not null"`
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
	Id        uint           `json:"id" gorm:"primaryKey"`
	UserId    uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// relationship
	CartItem []CartItem `json:"cart_item"`
}

type CartItem struct {
	Id        uint           `json:"id" gorm:"primaryKey"`
	CartId    uint           `json:"cart_id" gorm:"not null"`
	ProductId uint           `json:"product_id" gorm:"not null"`
	Quantity  uint           `json:"quantity" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// relationship
	Cart    Cart    `json:"-"`
	Product Product `json:"product"`
}

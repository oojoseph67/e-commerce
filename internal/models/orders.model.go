package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID           string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderNumber  string         `json:"order_number" gorm:"uniqueIndex;not null;size:50"`
	CustomerName string         `json:"customer_name" gorm:"not null"`
	UserID       string         `json:"user_id" gorm:"not null;type:uuid"`
	Status       OrderStatus    `json:"status" gorm:"default:pending"`
	TotalAmount  float64        `json:"total_amount" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

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
	Subtotal  float64        `json:"subtotal" gorm:"column:subtotal;not null"`
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
	CartItems []CartItem `json:"cart_items"`
}

type CartItem struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CartID    string         `json:"cart_id" gorm:"not null;type:uuid"`
	ProductID string         `json:"product_id" gorm:"not null;type:uuid"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// relationship
	Cart    Cart    `json:"-"`
	Product Product `json:"product"`
}

// type CartItem struct {
//     ProductID string
//     Prod      Product `gorm:"foreignKey:ProductID"` // field name is "Prod", not "Product"
// }

// // This will FAIL:
// db.Preload("CartItems.Product") // ❌ no field named "Product" on CartItem

// // This is correct:
// db.Preload("CartItems.Prod")    // ✅ matches the struct field name

// Convention (automatic):
// type CartItem struct {
//     ProductID string   // FK column: convention is <Assoc> + ID
//     Product   Product  // GORM auto-links ProductID → Product.ID
// }

// Explicit (when conventions don't match):
// type CartItem struct {
//     ProdID  string  `gorm:"column:item_product_id"`
//     Product Product `gorm:"foreignKey:ProdID;references:ID"`
// }

// type CartItem struct {
//     CartID  string  `gorm:"not null;type:uuid"`  // ← GORM sees this column
//     ProductID string `gorm:"not null;type:uuid"`  // ← GORM sees this column

//     Cart    Cart    `gorm:"foreignKey:CartID"`     // ← GORM links CartID → Cart.ID
//     Product Product `gorm:"foreignKey:ProductID"`  // ← GORM links ProductID → Product.ID
// }

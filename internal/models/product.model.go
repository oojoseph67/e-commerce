package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// relationships
	Products []Product `json:"-"`
}

type Product struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CategoryID  string         `json:"category_id" gorm:"not null;type:uuid"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Price       float64        `json:"price" gorm:"not null"`
	Stock       int            `json:"stock" gorm:"default:0"`
	SKU         string         `json:"sku" gorm:"uniqueIndex;not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// relationships
	Category   Category       `json:"category"`
	Images     []ProductImage `json:"images"`
	OrderItems []OrderItem    `json:"-"`
	CartItems  []CartItem     `json:"-"`
}

type ProductImage struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID string         `json:"product_id" gorm:"not null;type:uuid"`
	URL       string         `json:"url" gorm:"not null"`
	AltText   string         `json:"alt_text"`
	IsPrimary bool           `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// relationships
	Product Product `json:"-"`
}

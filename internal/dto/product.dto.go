package dto

// ProductImageUpload represents the uploaded product image data
type ProductImageUpload struct {
	URL     string `example:"https://storage.example.com/products/image.jpg"`
	AltText string `example:"Product image"`
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required" example:"Electronics"`
	Description string `json:"description" binding:"required" example:"Electronic devices and gadgets"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name" example:"Electronics"`
	Description string `json:"description" example:"Electronic devices and gadgets"`
	IsActive    *bool  `json:"is_active" example:"true"`
}

// CategoryResponse represents the category data in responses
type CategoryResponse struct {
	ID          string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name" example:"Electronics"`
	Description string `json:"description" example:"Electronic devices and gadgets"`
	IsActive    bool   `json:"is_active" example:"true"`
}

// GetCategoriesResponse represents the response for getting categories
type GetCategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int                `json:"total" example:"10"`
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	CategoryID  string  `json:"category_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string  `json:"name" binding:"required" example:"Wireless Headphones"`
	Description string  `json:"description" binding:"required" example:"Premium wireless headphones with noise cancellation"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"99.99"`
	Stock       int     `json:"stock" binding:"required,min=0" example:"100"`
	SKU         string  `json:"sku" binding:"required" example:"WH-1000XM4"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	CategoryID  string  `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string  `json:"name" example:"Wireless Headphones"`
	Description string  `json:"description" example:"Premium wireless headphones with noise cancellation"`
	Price       float64 `json:"price" binding:"omitempty,gt=0" example:"149.99"`
	Stock       int     `json:"stock" binding:"omitempty,min=0" example:"50"`
	IsActive    *bool   `json:"is_active" example:"true"`
}

// ProductResponse represents the product data in responses
type ProductResponse struct {
	ID          string                 `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CategoryID  string                 `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string                 `json:"name" example:"Wireless Headphones"`
	Description string                 `json:"description" example:"Premium wireless headphones with noise cancellation"`
	Price       float64                `json:"price" example:"99.99"`
	Stock       int                    `json:"stock" example:"100"`
	SKU         string                 `json:"sku" example:"WH-1000XM4"`
	IsActive    bool                   `json:"is_active" example:"true"`
	Category    CategoryResponse       `json:"category"`
	Images      []ProductImageResponse `json:"images"`
	Meta        ResponseMeta           `json:"meta"`
}

// ProductImageResponse represents the product image data in responses
type ProductImageResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	URL       string `json:"url" example:"https://storage.example.com/products/image.jpg"`
	AltText   string `json:"alt_text" example:"Product image"`
	IsPrimary bool   `json:"is_primary" example:"true"`
}

// ResponseMeta represents metadata for responses
type ResponseMeta struct {
	CreatedAt string `json:"created_at,omitempty" example:"2026-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at,omitempty" example:"2026-01-01T00:00:00Z"`
}

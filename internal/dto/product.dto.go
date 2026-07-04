package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type GetCategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int                `json:"total"`
}

type CreateProductRequest struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	SKU         string  `json:"sku" binding:"required"`
}

type UpdateProductRequest struct {
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,min=0"`
	IsActive    *bool   `json:"is_active"`
}

type ProductResponse struct {
	ID          string                 `json:"id"`
	CategoryID  string                 `json:"category_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Price       float64                `json:"price"`
	Stock       int                    `json:"stock"`
	SKU         string                 `json:"sku"`
	IsActive    bool                   `json:"is_active"`
	Category    CategoryResponse       `json:"category"`
	Images      []ProductImageResponse `json:"images"`
}

type ProductImageResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary"`
}

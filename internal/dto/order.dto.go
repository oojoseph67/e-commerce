package dto

type AddToCartRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quatity   int  `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quatity int `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	CartItems []CartItemResponse `json:"cart_item"`
	Total     float64            `json:"total"`
}

type CartItemResponse struct {
	ID       string          `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Subtotal float64         `json:"subtotal"`
}

type OrderResponse struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	OrderItems  []OrderItemResponse `json:"order_items"`
	CreatedAt   string              `json:"created_at"`
}

type OrderItemResponse struct {
	ID       string          `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Price    float64         `json:"price"`
}

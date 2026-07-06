package dto

// AddToCartRequest represents the request body for adding an item to cart
type AddToCartRequest struct {
	ProductID string `json:"product_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Quatity   int    `json:"quantity" binding:"required,min=1" example:"2"`
}

// UpdateCartItemRequest represents the request body for updating cart item quantity
type UpdateCartItemRequest struct {
	Quatity int `json:"quantity" binding:"required,min=1" example:"3"`
}

// CartResponse represents the cart data in responses
type CartResponse struct {
	ID             string             `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID         string             `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CartItems      []CartItemResponse `json:"cart_item"`
	Total          float64            `json:"total" example:"199.98"`
	TotalCartItems float64            `json:"total_cart_items" example:"2"`
}

// CartItemResponse represents a single cart item in responses
type CartItemResponse struct {
	ID       string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity" example:"2"`
	Subtotal float64         `json:"subtotal" example:"199.98"`
}

// OrderResponse represents the order data in responses
type OrderResponse struct {
	ID           string              `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderNumber  string              `json:"order_number" example:"ORD-20260101-001"`
	CustomerName string              `json:"customer_name" example:"John Doe"`
	UserID       string              `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status       string              `json:"status" example:"pending"`
	TotalAmount  float64             `json:"total_amount" example:"199.98"`
	OrderItems   []OrderItemResponse `json:"order_items"`
	User         UserResponse        `json:"user,omitzero"`
	CreatedAt    string              `json:"created_at" example:"2026-01-01T00:00:00Z"`
}

// OrdersResponse represents the response for getting multiple orders
type OrdersResponse struct {
	Orders []*OrderResponse `json:"orders"`
	User   *UserResponse    `json:"user"`
}

// OrderItemResponse represents a single order item in responses
type OrderItemResponse struct {
	ID       string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity" example:"2"`
	Subtotal float64         `json:"subtotal" example:"199.98"`
}

package handlers

import "github.com/oojoseph67/ecommerce/internal/services"

// type Handler struct {
// 	service *services.Services
// }

// func NewHandler(service *services.Services) *Handler {
// 	return &Handler{service: service}
// }

// so instead of using struct for typing we use interface to make sure the required methods are passed down correctly as well as there parameters and return types

type Handler struct {
	authService     services.AuthServicer
	userService     services.UserServicer
	productService  services.ProductServicer
	categoryService services.CategoryServicer
}

func NewHandler(
	auth services.AuthServicer,
	user services.UserServicer,
	product services.ProductServicer,
	category services.CategoryServicer,
) *Handler {
	return &Handler{
		authService:     auth,
		userService:     user,
		productService:  product,
		categoryService: category,
	}
}

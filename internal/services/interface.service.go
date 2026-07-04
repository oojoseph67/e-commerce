package services

import (
	"mime/multipart"

	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

// CONTAINS ALL THE EXPORTED METHODS FROM THE AUTHSERVICE, PRODUCTSERVICE, USERSERVICE, CATEGORYSERVICE
// commentforme: we are declaring the method and what we expect and return

type AuthServicer interface {
	AdminSignup(req *dto.AdminSignupRequest) (*dto.AuthResponse, error)
	Signup(req *dto.SignupRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	AdminLogin(req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(req *dto.LogoutRequest) error
}

type UserServicer interface {
	GetUserProfile(userId string) (*dto.UserResponse, error)
	UpdateProfile(userId string, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
}

type ProductServicer interface {
	CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProducts(page, limit int, category string) ([]dto.ProductResponse, *responses.PaginationMeta, error)
	GetProduct(id string) (*dto.ProductResponse, error)
	UpdateProduct(id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id string) error
	AddProductImages(productId string, images []dto.ProductImageUpload) error
	DeleteProductImage(imageId string) error
}

type CategoryServicer interface {
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories() (*dto.GetCategoriesResponse, error)
	UpdateCategory(id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id string) error
}

type UploadServicer interface {
	UploadProductImage(productId string, file *multipart.FileHeader) (url, altText string, err error)
}

type CartServicer interface {
	GetCart(userId string) (*dto.CartResponse, error)
	AddToCart(req *dto.AddToCartRequest, userId string) (*dto.CartResponse, error)
	UpdateCartItem(req *dto.UpdateCartItemRequest, cartItemId, userId string) (*dto.CartResponse, error)
	RemoveCartItem(cartItemId, userId string) error
}

type OrderServicer interface {
	CreateOrder(userId string) (*dto.OrderResponse, error)
	GetOrders(userId string, page, limit int) (*dto.OrdersResponse, *responses.PaginationMeta, error)
	GetOrder(userId, orderId string) (*dto.OrderResponse, error)
}

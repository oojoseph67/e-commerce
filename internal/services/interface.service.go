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
	GetProducts(page, limit int) ([]dto.ProductResponse, *responses.PaginationMeta, error)
	GetProduct(id string) (*dto.ProductResponse, error)
	UpdateProduct(id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id string) error
	AddProductImage(productId string, url, altText string) error
	DeleteProductImage(imageId string) error
}

type CategoryServicer interface {
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories() ([]dto.CategoryResponse, error)
	UpdateCategory(id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id string) error
}

type UploadServicer interface {
	UploadProductImage(productId string, file *multipart.FileHeader) (url, altText string, err error)
}

type CartServicer interface {
}

package services

import (
	"errors"

	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/models"
)

func (s *CartService) GetCart(userId string) (*dto.CartResponse, error) {
	// check if user exists
	var userModel models.User

	if err := s.db.Where("id = ? AND is_active = ?", userId, true).First(&userModel).Error; err != nil {
		s.internalLogger("cart").Warn().Str("userId", userId).Msg("user not found or is not active")
		return nil, errors.New("user not found or is not active")
	}

	// declare model
	var cartModel models.Cart
	err := s.db.Preload("CartItems.Product", "is_active = ?", true).
		Preload("CartItems.Product.Category").
		Preload("CartItems.Product.Images").
		Where("user_id = ?", userModel.ID).
		First(&cartModel).Error

	if err != nil {
		// create cart
		return s.createCart(userModel.ID)
	}

	return mapCartResponse(&cartModel)
}

func (s *CartService) AddToCart(req *dto.AddToCartRequest, userId string) (*dto.CartResponse, error) {
	// check that user exists
	if _, err := getUser(s.db, userId, *s.internalLogger("add_to_cart:check_user")); err != nil {
		return nil, err
	}

	// check that product exists
	product, err := getProductById(s.db, req.ProductID, *s.internalLogger("add_to_cart:check_product"))
	if err != nil {
		return nil, err
	}

	// check stock is more than enough
	if product.Stock < req.Quatity {
		return nil, errors.New("insufficient stock")
	}

	// get cart
	cart, err := s.GetCart(userId)
	if err != nil {
		return nil, err
	}

	var cartItemModel models.CartItem

	// cart item validation
	if err := s.db.Where("cart_id = ? AND product_id = ?", cart.ID, product.ID).First(&cartItemModel).Error; err != nil {

		// create cart item
		cartItemModel = models.CartItem{
			CartID:    cart.ID,
			ProductID: product.ID,
			Quantity:  req.Quatity,
		}

		if err := s.db.Create(&cartItemModel).Error; err != nil {
			s.internalLogger("add_to_cart:add_new_item").Error().Str("product_id", req.ProductID).Int("quantity", req.Quatity).Err(err).Msg("failed to add product to cart")
			return nil, errors.New("failed to add item to cart")
		}
	} else {
		// update exiting cart item if item already exists
		cartItemModel.Quantity += req.Quatity

		// check stock is more than enough
		if product.Stock < req.Quatity {
			return nil, errors.New("insufficient stock")
		}

		err := s.db.Model(cartItemModel).Updates(map[string]interface{}{
			"quantity": cartItemModel.Quantity,
		}).Error

		if err != nil {
			s.internalLogger("add_to_cart:add_new_item").Error().Str("product_id", req.ProductID).Int("old_quantity", req.Quatity).Int("new_quantity", cartItemModel.Quantity).Err(err).Msg("failed to update item cart")
			return nil, nil
		}
	}

	return s.GetCart(userId)
}

func (s *CartService) UpdateCartItem(req *dto.UpdateCartItemRequest, cartItemId, userId string) (*dto.CartResponse, error) {
	// get cart
	cart, err := s.GetCart(userId)
	if err != nil {
		return nil, err
	}

	// check cart-item exists
	var cartItemModel models.CartItem

	if err := s.db.Where("cart_id = ? AND id = ?", cart.ID, cartItemId).First(&cartItemModel).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	// check that product exists
	product, err := getProductById(s.db, cartItemModel.ProductID, *s.internalLogger("update_cart_id:check_product"))
	if err != nil {
		return nil, err
	}

	// check stock is more than enough
	if product.Stock < req.Quatity {
		return nil, errors.New("insufficient stock")
	}

	// update exiting cart item if item already exists
	cartItemModel.Quantity += req.Quatity

	err = s.db.Model(cartItemModel).Updates(map[string]interface{}{
		"quantity": cartItemModel.Quantity,
	}).Error

	if err != nil {
		s.internalLogger("add_to_cart:add_new_item").Error().Str("product_id", cartItemModel.ProductID).Int("old_quantity", req.Quatity).Int("new_quantity", cartItemModel.Quantity).Err(err).Msg("failed to update item cart")
		return nil, nil
	}

	return s.GetCart(userId)

}

func (s *CartService) RemoveCartItem(cartItemId, userId string) error {
	// get cart
	cart, err := s.GetCart(userId)
	if err != nil {
		return err
	}

	// check cart-item exists
	var cartItemModel models.CartItem
	if err := s.db.Where("id = ? AND cart_id = ?", cartItemId, cart.ID).Delete(&cartItemModel).Error; err != nil {
		s.internalLogger("cart:remove_cart_id").Warn().Str("cart_item_id", cartItemId).Str("user_id", userId).Err(err).Msg("error removing cart-item")
		return errors.New("error removing item from cart")
	}

	return nil
}

func (s *CartService) createCart(userId string) (*dto.CartResponse, error) {
	// map model
	userCart := models.Cart{
		UserID: userId,
	}

	// save model
	if err := s.db.Create(&userCart).Error; err != nil {
		s.internalLogger("cart:create").Error().Err(err).Msg("error creating cart")
		return nil, err
	}

	return nil, nil
}

func mapCartResponse(cartModel *models.Cart) (*dto.CartResponse, error) {

	var total float64
	cartItems := make([]dto.CartItemResponse, len(cartModel.CartItems))

	for i := range cartModel.CartItems {
		subtotal := float64(cartModel.CartItems[i].Quantity) * cartModel.CartItems[i].Product.Price
		total += subtotal
		product := convertToProductResponse(&cartModel.CartItems[i].Product)

		cartItems[i] = dto.CartItemResponse{
			ID:       cartModel.CartItems[i].ID,
			Quantity: int(cartModel.CartItems[i].Quantity),
			Subtotal: subtotal,
			Product:  *product,
		}
	}

	cart := &dto.CartResponse{
		ID:             cartModel.ID,
		UserID:         cartModel.UserID,
		CartItems:      cartItems,
		Total:          total,
		TotalCartItems: float64(len(cartModel.CartItems)),
	}

	return cart, nil
}

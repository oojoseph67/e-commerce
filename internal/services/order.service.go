package services

import (
	"errors"

	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *OrderService) CreateOrder(userId string) (*dto.OrderResponse, error) {

	var orderResponse *dto.OrderResponse

	var userModel models.User
	var cartModel models.Cart
	var cartItemsModel []models.CartItem
	var orderModel models.Order

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// check userId
		if err := tx.Where("id = ? AND is_active = ?", userId, true).First(&userModel).Error; err != nil {
			s.internalLogger("order:create_order").Warn().Str("user_id", userId).Err(err).Msg("user not found")
			return err
		}

		// get cart
		if err := tx.Where("user_id = ?", userModel.ID).First(&cartModel).Error; err != nil {
			s.internalLogger("order:create_order").Warn().Str("user_id", userId).Err(err).Msg("couldnt find user cart")
			return err
		}

		// get cart items
		if err := tx.Where("cart_id = ?", cartModel.ID).Find(&cartItemsModel).Error; err != nil {
			s.internalLogger("order:create_order").Warn().Str("cart_id", cartModel.ID).Err(err).Msg("couldnt find cart items")
			return err
		}

		// check that there are items in the cart
		if len(cartItemsModel) == 0 {
			return errors.New("cant create empty cart")
		}

		// validate product stock, update stock, and build order items in one pass
		var total float64
		var orderItemsToCreate []models.OrderItem

		for _, cartItem := range cartItemsModel {
			var productModel models.Product

			// locking the product query
			err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Where("id = ?", cartItem.Product.ID).First(&productModel).Error
			if err != nil {
				s.internalLogger("order:create_order").Warn().Str("product_id", cartItem.Product.ID).Err(err).Msg("couldnt cart item product")
				return err
			}

			subtotal := float64(cartItem.Quantity) * cartItem.Product.Price
			total += subtotal

			// subtract quantity from product stock
			if productModel.Stock >= cartItem.Quantity {
				productModel.Stock = productModel.Stock - cartItem.Quantity
			} else {
				err := errors.New("insufficient stock")
				s.internalLogger("order:create_order").Warn().Str("product_id", cartItem.Product.ID).Err(err).Msg("insufficient stocke")
				return err
			}

			// update the stock
			if err := tx.Model(&productModel).Update("stock", productModel.Stock).Error; err != nil {
				s.internalLogger("order:create_order").Warn().Str("product_id", cartItem.Product.ID).Err(err).Msg("failed tp update product stock")
				return err
			}

			// build order item for batch creation (order id assigned after order creation)
			orderItemsToCreate = append(orderItemsToCreate, models.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Subtotal:  subtotal,
			})
		}

		// map order
		orderModel = models.Order{
			UserID:      userModel.ID,
			TotalAmount: total,
		}

		// create order
		if err := tx.Create(&orderModel).Error; err != nil {
			s.internalLogger("order:create_order").Warn().Err(err).Msg("error creating order")
			return err
		}

		// assign order id to each order item
		for i := range orderItemsToCreate {
			orderItemsToCreate[i].OrderID = orderModel.ID
		}

		// create order items in batch
		if err := tx.CreateInBatches(orderItemsToCreate, 100).Error; err != nil {
			s.internalLogger("order:create_order").Warn().Err(err).Msg("error creating order items")
			return err
		}

		// clear cart after order has been created
		if err := tx.Where("cart_id = ?", cartModel.ID).Delete(&models.CartItem{}).Error; err != nil {
			s.internalLogger("order:create_order").Warn().Err(err).Msg("failed to clear cart")
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := s.db.Where("id = ? AND user_id = ?", orderModel.ID, userId).First(&orderModel).Error; err != nil {
		s.internalLogger("order:create_order").Warn().Str("user_id", userId).Str("order_id", orderModel.ID).Err(err).Msg("couldnt re-query order ")
		return nil, err
	}

	orderItems := make([]dto.OrderItemResponse, len(orderModel.OrderItems))

	for i := range orderModel.OrderItems {
		orderItems[i] = dto.OrderItemResponse{
			ID:       orderModel.OrderItems[i].ID,
			Product:  *convertToProductResponse(&orderModel.OrderItems[i].Product),
			Quantity: orderModel.OrderItems[i].Quantity,
			Subtotal: orderModel.OrderItems[i].Subtotal,
		}
	}

	orderResponse = &dto.OrderResponse{
		ID:          orderModel.ID,
		UserID:      userModel.ID,
		Status:      string(orderModel.Status),
		TotalAmount: orderModel.TotalAmount,
		OrderItems:  orderItems,
	}

	return orderResponse, nil
}

package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/models"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
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
		if err := tx.Preload("Product").Where("cart_id = ?", cartModel.ID).Find(&cartItemsModel).Error; err != nil {
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
				productModel.Stock -= cartItem.Quantity
			} else {
				err := errors.New("insufficient stock")
				s.internalLogger("order:create_order").Warn().Str("product_id", cartItem.Product.ID).Err(err).Msg("insufficient stock")
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

		// generate order number
		orderNumber := fmt.Sprintf("ORD-%d", time.Now().UnixNano()%1000000)
		customerName := fmt.Sprintf("%s %s", userModel.FirstName, userModel.LastName)

		// map order
		orderModel = models.Order{
			OrderNumber:  orderNumber,
			CustomerName: customerName,
			UserID:       userModel.ID,
			TotalAmount:  total,
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

	if err := s.db.Preload("OrderItems.Product.Category").Where("id = ? AND user_id = ?", orderModel.ID, userId).First(&orderModel).Error; err != nil {
		s.internalLogger("order:create_order").Warn().Str("user_id", userId).Str("order_id", orderModel.ID).Err(err).Msg("couldnt re-query order ")
		return nil, err
	}

	orderResponse = covertToOrderResponse(&orderModel)

	return orderResponse, nil
}

func (s *OrderService) GetOrders(userId string, page, limit int) ([]*dto.OrderResponse, *responses.PaginationMeta, error) {
	offset := (page - 1) * limit

	var total int64
	var userModel models.User
	var ordersModel []models.Order

	// check user exists
	if err := s.db.Where("id = ? AND is_active = ?", userId, true).First(&userModel).Error; err != nil {
		s.internalLogger("order:get_order").Warn().Str("user_id", userId).Err(err).Msg("user not found")
		return []*dto.OrderResponse{}, nil, err
	}

	// check order exists
	s.db.Model(&models.Order{}).Where("user_id = ?", userId).Count(&total)

	err := s.db.Preload("OrderItems.Product.Category").
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&ordersModel).Error

	if err != nil {
		s.internalLogger("order:get_order").Warn().Str("user_id", userId).Err(err).Msg("couldnt get user orders")
		return nil, nil, nil
	}

	var ordersResponse []*dto.OrderResponse

	for i := range ordersModel {
		order := ordersModel[i]
		response := covertToOrderResponse(&order)
		ordersResponse = append(ordersResponse, response)
	}

	totalPages := int(total+int64(limit)-1) / int(limit)
	meta := &responses.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return ordersResponse, meta, nil
}

func (s *OrderService) GetOrder(userId, orderId string) (*dto.OrderResponse, error) {
	var userModel models.User
	var orderModel models.Order

	// check user exists
	if err := s.db.Where("id = ? AND is_active = ?", userId, true).First(&userModel).Error; err != nil {
		s.internalLogger("order:get_order").Warn().Str("user_id", userId).Err(err).Msg("user not found")
		return nil, err
	}

	// check order exists
	if err := s.db.Preload("OrderItems.Product.Category").Where("id = ? AND user_id = ?", orderId, userId).First(&orderModel).Error; err != nil {
		s.internalLogger("order:get_order").Warn().Str("user_id", userId).Str("order_id", orderId).Err(err).Msg("order not found for user")
		return nil, err
	}

	orderResponse := covertToOrderResponse(&orderModel)

	return orderResponse, nil
}

func covertToOrderResponse(order *models.Order) *dto.OrderResponse {
	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))

	for i := range order.OrderItems {
		orderItems[i] = dto.OrderItemResponse{
			ID:       order.OrderItems[i].ID,
			Product:  *convertToProductResponse(&order.OrderItems[i].Product),
			Quantity: order.OrderItems[i].Quantity,
			Subtotal: order.OrderItems[i].Subtotal,
		}
	}

	orderResponse := &dto.OrderResponse{
		ID:           order.ID,
		OrderNumber:  order.OrderNumber,
		CustomerName: order.CustomerName,
		UserID:       order.UserID,
		Status:       string(order.Status),
		TotalAmount:  order.TotalAmount,
		OrderItems:   orderItems,
		CreatedAt:    order.CreatedAt.Format(defaultDateFormat),
	}
	return orderResponse
}

package services

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/repositories"
	"fmt"
)

// OrderService handles order-related business logic
type OrderService struct {
	orderRepo    *repositories.OrderRepository
	carRepo      *repositories.CarRepository
	emailService *EmailService
	userRepo     *repositories.UserRepository
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo *repositories.OrderRepository,
	carRepo *repositories.CarRepository,
	emailService *EmailService,
	userRepo *repositories.UserRepository,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		carRepo:      carRepo,
		emailService: emailService,
		userRepo:     userRepo,
	}
}

// AddToCart adds a car to user's cart
func (s *OrderService) AddToCart(userID, carID uint) error {
	// Check if car exists and is available
	car, err := s.carRepo.GetByID(carID)
	if err != nil {
		return fmt.Errorf("car not found")
	}

	if !car.IsAvailable {
		return fmt.Errorf("car is not available")
	}

	// Add to cart
	if err := s.orderRepo.AddToCart(userID, carID); err != nil {
		return fmt.Errorf("failed to add to cart: %w", err)
	}

	return nil
}

// RemoveFromCart removes a car from user's cart
func (s *OrderService) RemoveFromCart(userID, carID uint) error {
	if err := s.orderRepo.RemoveFromCart(userID, carID); err != nil {
		return fmt.Errorf("failed to remove from cart: %w", err)
	}

	return nil
}

// GetCart retrieves user's cart
func (s *OrderService) GetCart(userID uint) (*models.Cart, error) {
	cart, err := s.orderRepo.GetCart(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return cart, nil
}

// CreateOrder creates an order from user's cart
func (s *OrderService) CreateOrder(userID uint) (*models.Order, error) {
	// Get cart
	cart, err := s.orderRepo.GetCart(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Verify all cars are still available
	for _, item := range cart.Items {
		car, err := s.carRepo.GetByID(item.CarID)
		if err != nil {
			return nil, fmt.Errorf("car %d not found", item.CarID)
		}

		if !car.IsAvailable {
			return nil, fmt.Errorf("car %s %s is no longer available", car.Make, car.Model)
		}
	}

	// Create order
	order := &models.Order{
		UserID:      userID,
		TotalAmount: cart.TotalAmount,
		Status:      "pending",
		Items:       make([]models.OrderItem, len(cart.Items)),
	}

	// Convert cart items to order items
	for i, cartItem := range cart.Items {
		order.Items[i] = models.OrderItem{
			CarID: cartItem.CarID,
			Price: cartItem.Car.Price,
		}
	}

	// Save order
	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Mark cars as unavailable
	for _, item := range order.Items {
		car, _ := s.carRepo.GetByID(item.CarID)
		if car != nil {
			car.IsAvailable = false
			_ = s.carRepo.Update(car)
		}
	}

	// Clear cart
	if err := s.orderRepo.ClearCart(userID); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Send order confirmation email (non-blocking)
	go func() {
		user, err := s.userRepo.GetByID(userID)
		if err == nil {
			_ = s.emailService.SendOrderConfirmation(user.Email, user.FirstName, order.ID, order.TotalAmount)
		}
	}()

	return order, nil
}

// GetOrderHistory retrieves user's order history
func (s *OrderService) GetOrderHistory(userID uint) ([]*models.Order, error) {
	orders, err := s.orderRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}

	return orders, nil
}

// GetOrderByID retrieves an order by ID
func (s *OrderService) GetOrderByID(orderID uint) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(orderID uint, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"completed":  true,
		"cancelled":  true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status")
	}

	// Get order
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	// Update status
	order.Status = status

	if err := s.orderRepo.Update(order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// If order is cancelled, mark cars as available again
	if status == "cancelled" {
		for _, item := range order.Items {
			car, err := s.carRepo.GetByID(item.CarID)
			if err == nil {
				car.IsAvailable = true
				_ = s.carRepo.Update(car)
			}
		}
	}

	return nil
}

// GetAllOrders retrieves all orders (admin function)
func (s *OrderService) GetAllOrders(limit, offset int) ([]*models.Order, error) {
	if limit == 0 {
		limit = 50
	}

	orders, err := s.orderRepo.GetAllOrders(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, nil
}

// GetOrderStatistics retrieves order statistics for a time period
func (s *OrderService) GetOrderStatistics(period models.TimePeriod) (*models.OrderStats, error) {
	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return stats, nil
}

// ClearCart clears user's cart
func (s *OrderService) ClearCart(userID uint) error {
	if err := s.orderRepo.ClearCart(userID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

// GetCartItemCount returns the number of items in user's cart
func (s *OrderService) GetCartItemCount(userID uint) (int, error) {
	cart, err := s.orderRepo.GetCart(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get cart: %w", err)
	}

	return cart.ItemCount, nil
}

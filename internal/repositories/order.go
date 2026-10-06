package repositories

import (
	"car-dealership-system/internal/models"
	"database/sql"
	"fmt"
)

// OrderRepository handles order and cart data operations
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order with items
func (r *OrderRepository) Create(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert order
	query := `
		INSERT INTO orders (user_id, total_amount, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(query, order.UserID, order.TotalAmount, order.Status).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (order_id, car_id, price)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	for i := range order.Items {
		item := &order.Items[i]
		item.OrderID = order.ID

		err = tx.QueryRow(itemQuery, item.OrderID, item.CarID, item.Price).Scan(
			&item.ID,
			&item.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByUserID retrieves all orders for a user
func (r *OrderRepository) GetByUserID(userID uint) ([]*models.Order, error) {
	query := `
		SELECT id, user_id, total_amount, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Load order items
		items, err := r.getOrderItems(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, &order)
	}

	return orders, nil
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(id uint) (*models.Order, error) {
	query := `
		SELECT id, user_id, total_amount, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order
	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Load order items
	items, err := r.getOrderItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}

// getOrderItems retrieves items for an order
func (r *OrderRepository) getOrderItems(orderID uint) ([]models.OrderItem, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.car_id, oi.price, oi.created_at,
		       c.make, c.model, c.year, c.color, c.image_url
		FROM order_items oi
		JOIN cars c ON oi.car_id = c.id
		WHERE oi.order_id = $1
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var car models.Car

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.CarID,
			&item.Price,
			&item.CreatedAt,
			&car.Make,
			&car.Model,
			&car.Year,
			&car.Color,
			&car.ImageURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		car.ID = item.CarID
		item.Car = &car
		items = append(items, item)
	}

	return items, nil
}

// Update updates an order
func (r *OrderRepository) Update(order *models.Order) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, order.Status, order.ID).Scan(&order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// AddToCart adds a car to user's cart
func (r *OrderRepository) AddToCart(userID, carID uint) error {
	query := `
		INSERT INTO cart_items (user_id, car_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, car_id) DO NOTHING
	`

	_, err := r.db.Exec(query, userID, carID)
	if err != nil {
		return fmt.Errorf("failed to add to cart: %w", err)
	}

	return nil
}

// RemoveFromCart removes a car from user's cart
func (r *OrderRepository) RemoveFromCart(userID, carID uint) error {
	query := `DELETE FROM cart_items WHERE user_id = $1 AND car_id = $2`

	result, err := r.db.Exec(query, userID, carID)
	if err != nil {
		return fmt.Errorf("failed to remove from cart: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("item not in cart")
	}

	return nil
}

// GetCart retrieves user's cart
func (r *OrderRepository) GetCart(userID uint) (*models.Cart, error) {
	query := `
		SELECT ci.id, ci.user_id, ci.car_id, ci.created_at,
		       c.make, c.model, c.year, c.price, c.color, c.image_url, c.is_available
		FROM cart_items ci
		JOIN cars c ON ci.car_id = c.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	defer rows.Close()

	cart := &models.Cart{
		Items: []models.CartItem{},
	}

	for rows.Next() {
		var item models.CartItem
		var car models.Car

		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.CarID,
			&item.CreatedAt,
			&car.Make,
			&car.Model,
			&car.Year,
			&car.Price,
			&car.Color,
			&car.ImageURL,
			&car.IsAvailable,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}

		car.ID = item.CarID
		item.Car = &car
		cart.Items = append(cart.Items, item)
		cart.TotalAmount += car.Price
		cart.ItemCount++
	}

	return cart, nil
}

// ClearCart removes all items from user's cart
func (r *OrderRepository) ClearCart(userID uint) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

// GetStatistics retrieves order statistics for a time period
func (r *OrderRepository) GetStatistics(period models.TimePeriod) (*models.OrderStats, error) {
	stats := &models.OrderStats{
		OrdersByStatus: make(map[string]int),
		OrdersByMonth:  make(map[string]models.OrderMetric),
	}

	// Total orders and revenue
	query := `
		SELECT COUNT(*), COALESCE(SUM(total_amount), 0)
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
	`

	err := r.db.QueryRow(query, period.StartDate, period.EndDate).Scan(
		&stats.TotalOrders,
		&stats.TotalRevenue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get total stats: %w", err)
	}

	if stats.TotalOrders > 0 {
		stats.AverageOrder = stats.TotalRevenue / float64(stats.TotalOrders)
	}

	// Orders by status
	statusQuery := `
		SELECT status, COUNT(*)
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
		GROUP BY status
	`

	rows, err := r.db.Query(statusQuery, period.StartDate, period.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get status stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.OrdersByStatus[status] = count
	}

	// Orders by month
	monthQuery := `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM') as month,
			COUNT(*),
			COALESCE(SUM(total_amount), 0)
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
		GROUP BY month
		ORDER BY month
	`

	rows, err = r.db.Query(monthQuery, period.StartDate, period.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get month stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var month string
		var metric models.OrderMetric
		if err := rows.Scan(&month, &metric.Count, &metric.Revenue); err != nil {
			return nil, err
		}
		stats.OrdersByMonth[month] = metric
	}

	// Top cars
	topCarsQuery := `
		SELECT 
			c.id, c.make, c.model,
			COUNT(oi.id) as order_count,
			COALESCE(SUM(oi.price), 0) as revenue
		FROM order_items oi
		JOIN cars c ON oi.car_id = c.id
		JOIN orders o ON oi.order_id = o.id
		WHERE o.created_at BETWEEN $1 AND $2
		GROUP BY c.id, c.make, c.model
		ORDER BY order_count DESC
		LIMIT 10
	`

	rows, err = r.db.Query(topCarsQuery, period.StartDate, period.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get top cars: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var carStat models.CarStat
		if err := rows.Scan(&carStat.CarID, &carStat.Make, &carStat.Model, &carStat.OrderCount, &carStat.Revenue); err != nil {
			return nil, err
		}
		stats.TopCars = append(stats.TopCars, carStat)
	}

	return stats, nil
}

// GetAllOrders retrieves all orders (admin function)
func (r *OrderRepository) GetAllOrders(limit, offset int) ([]*models.Order, error) {
	query := `
		SELECT o.id, o.user_id, o.total_amount, o.status, o.created_at, o.updated_at,
		       u.email, u.first_name, u.last_name
		FROM orders o
		JOIN users u ON o.user_id = u.id
		ORDER BY o.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		var user models.User

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
			&user.Email,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		user.ID = order.UserID
		order.User = &user

		// Load order items
		items, err := r.getOrderItems(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, &order)
	}

	return orders, nil
}

// GetOrderStatusBreakdown returns order counts grouped by status
func (r *OrderRepository) GetOrderStatusBreakdown() (map[string]int, error) {
	query := `
		SELECT status, COUNT(*)
		FROM orders
		GROUP BY status
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status breakdown: %w", err)
	}
	defer rows.Close()

	breakdown := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan status breakdown: %w", err)
		}
		breakdown[status] = count
	}

	// Ensure all statuses are present
	statuses := []string{"pending", "processing", "completed", "cancelled"}
	for _, status := range statuses {
		if _, exists := breakdown[status]; !exists {
			breakdown[status] = 0
		}
	}

	return breakdown, nil
}

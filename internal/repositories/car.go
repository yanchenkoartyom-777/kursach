package repositories

import (
	"car-dealership-system/internal/models"
	"database/sql"
	"fmt"
	"strings"
)

// CarRepository handles car data operations
type CarRepository struct {
	db *sql.DB
}

// NewCarRepository creates a new car repository
func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{db: db}
}

// Create creates a new car
func (r *CarRepository) Create(car *models.Car) error {
	query := `
		INSERT INTO cars (make, model, year, price, color, engine_type, transmission, mileage, description, image_url, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		car.Make,
		car.Model,
		car.Year,
		car.Price,
		car.Color,
		car.EngineType,
		car.Transmission,
		car.Mileage,
		car.Description,
		car.ImageURL,
		car.IsAvailable,
	).Scan(&car.ID, &car.CreatedAt, &car.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	return nil
}

// GetAll retrieves all cars with optional filters
func (r *CarRepository) GetAll(filters models.CarFilters) ([]*models.Car, error) {
	query := `
		SELECT id, make, model, year, price, color, engine_type, transmission, mileage, description, image_url, is_available, created_at, updated_at
		FROM cars
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	// Apply filters
	if filters.Make != "" {
		query += fmt.Sprintf(" AND LOWER(make) = LOWER($%d)", argCount)
		args = append(args, filters.Make)
		argCount++
	}

	if filters.Model != "" {
		query += fmt.Sprintf(" AND LOWER(model) = LOWER($%d)", argCount)
		args = append(args, filters.Model)
		argCount++
	}

	if filters.MinYear > 0 {
		query += fmt.Sprintf(" AND year >= $%d", argCount)
		args = append(args, filters.MinYear)
		argCount++
	}

	if filters.MaxYear > 0 {
		query += fmt.Sprintf(" AND year <= $%d", argCount)
		args = append(args, filters.MaxYear)
		argCount++
	}

	if filters.MinPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, filters.MinPrice)
		argCount++
	}

	if filters.MaxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, filters.MaxPrice)
		argCount++
	}

	if filters.EngineType != "" {
		query += fmt.Sprintf(" AND LOWER(engine_type) = LOWER($%d)", argCount)
		args = append(args, filters.EngineType)
		argCount++
	}

	if filters.Transmission != "" {
		query += fmt.Sprintf(" AND LOWER(transmission) = LOWER($%d)", argCount)
		args = append(args, filters.Transmission)
		argCount++
	}

	if filters.IsAvailable != nil {
		query += fmt.Sprintf(" AND is_available = $%d", argCount)
		args = append(args, *filters.IsAvailable)
		argCount++
	}

	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		switch filters.SortBy {
		case "price", "year", "created_at":
			sortBy = filters.SortBy
		}
	}

	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Apply pagination
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.Limit)
		argCount++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filters.Offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}
	defer rows.Close()

	var cars []*models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(
			&car.ID,
			&car.Make,
			&car.Model,
			&car.Year,
			&car.Price,
			&car.Color,
			&car.EngineType,
			&car.Transmission,
			&car.Mileage,
			&car.Description,
			&car.ImageURL,
			&car.IsAvailable,
			&car.CreatedAt,
			&car.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan car: %w", err)
		}
		cars = append(cars, &car)
	}

	return cars, nil
}

// GetByID retrieves a car by ID
func (r *CarRepository) GetByID(id uint) (*models.Car, error) {
	query := `
		SELECT id, make, model, year, price, color, engine_type, transmission, mileage, description, image_url, is_available, created_at, updated_at
		FROM cars
		WHERE id = $1
	`

	var car models.Car
	err := r.db.QueryRow(query, id).Scan(
		&car.ID,
		&car.Make,
		&car.Model,
		&car.Year,
		&car.Price,
		&car.Color,
		&car.EngineType,
		&car.Transmission,
		&car.Mileage,
		&car.Description,
		&car.ImageURL,
		&car.IsAvailable,
		&car.CreatedAt,
		&car.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("car not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	return &car, nil
}

// Update updates a car
func (r *CarRepository) Update(car *models.Car) error {
	query := `
		UPDATE cars
		SET make = $1, model = $2, year = $3, price = $4, color = $5, engine_type = $6, 
		    transmission = $7, mileage = $8, description = $9, image_url = $10, is_available = $11,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		car.Make,
		car.Model,
		car.Year,
		car.Price,
		car.Color,
		car.EngineType,
		car.Transmission,
		car.Mileage,
		car.Description,
		car.ImageURL,
		car.IsAvailable,
		car.ID,
	).Scan(&car.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}

	return nil
}

// Delete deletes a car
func (r *CarRepository) Delete(id uint) error {
	query := `DELETE FROM cars WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("car not found")
	}

	return nil
}

// Search searches for cars by query string
func (r *CarRepository) Search(query models.SearchQuery) ([]*models.Car, error) {
	sqlQuery := `
		SELECT id, make, model, year, price, color, engine_type, transmission, mileage, description, image_url, is_available, created_at, updated_at
		FROM cars
		WHERE (
			LOWER(make) LIKE LOWER($1) OR
			LOWER(model) LIKE LOWER($1) OR
			LOWER(description) LIKE LOWER($1)
		)
	`

	searchPattern := "%" + strings.ToLower(query.Query) + "%"
	args := []interface{}{searchPattern}
	argCount := 2

	// Apply additional filters
	if query.Filters.IsAvailable != nil {
		sqlQuery += fmt.Sprintf(" AND is_available = $%d", argCount)
		args = append(args, *query.Filters.IsAvailable)
		argCount++
	}

	if query.Filters.MinPrice > 0 {
		sqlQuery += fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, query.Filters.MinPrice)
		argCount++
	}

	if query.Filters.MaxPrice > 0 {
		sqlQuery += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, query.Filters.MaxPrice)
		argCount++
	}

	// Apply sorting
	sortBy := "created_at"
	if query.Filters.SortBy != "" {
		switch query.Filters.SortBy {
		case "price", "year", "created_at":
			sortBy = query.Filters.SortBy
		}
	}

	sortOrder := "DESC"
	if query.Filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	sqlQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	if query.Filters.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, query.Filters.Limit)
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search cars: %w", err)
	}
	defer rows.Close()

	var cars []*models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(
			&car.ID,
			&car.Make,
			&car.Model,
			&car.Year,
			&car.Price,
			&car.Color,
			&car.EngineType,
			&car.Transmission,
			&car.Mileage,
			&car.Description,
			&car.ImageURL,
			&car.IsAvailable,
			&car.CreatedAt,
			&car.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan car: %w", err)
		}
		cars = append(cars, &car)
	}

	return cars, nil
}

// Count returns the total number of cars matching filters
func (r *CarRepository) Count(filters models.CarFilters) (int, error) {
	query := `SELECT COUNT(*) FROM cars WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if filters.Make != "" {
		query += fmt.Sprintf(" AND LOWER(make) = LOWER($%d)", argCount)
		args = append(args, filters.Make)
		argCount++
	}

	if filters.IsAvailable != nil {
		query += fmt.Sprintf(" AND is_available = $%d", argCount)
		args = append(args, *filters.IsAvailable)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count cars: %w", err)
	}

	return count, nil
}

// GetInventoryStats returns car inventory statistics
func (r *CarRepository) GetInventoryStats() (*models.CarInventoryStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_cars,
			COUNT(*) FILTER (WHERE is_available = true) as available_cars,
			COUNT(*) FILTER (WHERE is_available = false) as unavailable_cars
		FROM cars
	`

	var stats models.CarInventoryStats
	err := r.db.QueryRow(query).Scan(
		&stats.TotalCars,
		&stats.AvailableCars,
		&stats.UnavailableCars,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get inventory stats: %w", err)
	}

	// Calculate sold cars (cars that are in completed orders)
	soldQuery := `
		SELECT COUNT(DISTINCT oi.car_id)
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE o.status = 'completed'
	`

	err = r.db.QueryRow(soldQuery).Scan(&stats.SoldCars)
	if err != nil {
		return nil, fmt.Errorf("failed to get sold cars count: %w", err)
	}

	return &stats, nil
}

// GetSoldCarsCount returns count of sold cars
func (r *CarRepository) GetSoldCarsCount() (int, error) {
	query := `
		SELECT COUNT(DISTINCT oi.car_id)
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE o.status = 'completed'
	`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get sold cars count: %w", err)
	}

	return count, nil
}

package services

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/repositories"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// CarService handles car-related business logic
type CarService struct {
	repo      *repositories.CarRepository
	validator *validator.Validate
}

// NewCarService creates a new car service
func NewCarService(repo *repositories.CarRepository) *CarService {
	return &CarService{
		repo:      repo,
		validator: validator.New(),
	}
}

// CreateCar creates a new car
func (s *CarService) CreateCar(car *models.Car) error {
	// Validate car data
	if err := s.validator.Struct(car); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Additional business validations
	if car.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	if car.Year < 1900 || car.Year > 2100 {
		return fmt.Errorf("invalid year")
	}

	// Create car
	if err := s.repo.Create(car); err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	return nil
}

// GetCars retrieves cars with optional filters
func (s *CarService) GetCars(filters models.CarFilters) ([]*models.Car, error) {
	// Set default limit if not specified
	if filters.Limit == 0 {
		filters.Limit = 50
	}

	cars, err := s.repo.GetAll(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}

	return cars, nil
}

// GetCarByID retrieves a car by ID
func (s *CarService) GetCarByID(id uint) (*models.Car, error) {
	car, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("car not found")
	}

	return car, nil
}

// UpdateCar updates a car
func (s *CarService) UpdateCar(car *models.Car) error {
	// Validate car data
	if err := s.validator.Struct(car); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Check if car exists
	_, err := s.repo.GetByID(car.ID)
	if err != nil {
		return fmt.Errorf("car not found")
	}

	// Update car
	if err := s.repo.Update(car); err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}

	return nil
}

// DeleteCar deletes a car
func (s *CarService) DeleteCar(id uint) error {
	// Check if car exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("car not found")
	}

	// Delete car
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	return nil
}

// SearchCars searches for cars by query string
func (s *CarService) SearchCars(queryStr string, filters models.CarFilters) ([]*models.Car, error) {
	if queryStr == "" {
		return s.GetCars(filters)
	}

	// Set default limit if not specified
	if filters.Limit == 0 {
		filters.Limit = 50
	}

	query := models.SearchQuery{
		Query:   queryStr,
		Filters: filters,
	}

	cars, err := s.repo.Search(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search cars: %w", err)
	}

	return cars, nil
}

// GetAvailableCars retrieves only available cars
func (s *CarService) GetAvailableCars(filters models.CarFilters) ([]*models.Car, error) {
	available := true
	filters.IsAvailable = &available

	return s.GetCars(filters)
}

// MarkAsUnavailable marks a car as unavailable
func (s *CarService) MarkAsUnavailable(id uint) error {
	car, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("car not found")
	}

	car.IsAvailable = false

	if err := s.repo.Update(car); err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}

	return nil
}

// MarkAsAvailable marks a car as available
func (s *CarService) MarkAsAvailable(id uint) error {
	car, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("car not found")
	}

	car.IsAvailable = true

	if err := s.repo.Update(car); err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}

	return nil
}

// GetCarCount returns the total number of cars matching filters
func (s *CarService) GetCarCount(filters models.CarFilters) (int, error) {
	count, err := s.repo.Count(filters)
	if err != nil {
		return 0, fmt.Errorf("failed to count cars: %w", err)
	}

	return count, nil
}

// GetMakes returns a list of unique car makes
func (s *CarService) GetMakes() ([]string, error) {
	// Get all cars
	cars, err := s.repo.GetAll(models.CarFilters{Limit: 1000})
	if err != nil {
		return nil, fmt.Errorf("failed to get makes: %w", err)
	}

	// Extract unique makes
	makesMap := make(map[string]bool)
	for _, car := range cars {
		makesMap[car.Make] = true
	}

	makes := make([]string, 0, len(makesMap))
	for make := range makesMap {
		makes = append(makes, make)
	}

	return makes, nil
}

// GetModels returns a list of models for a specific make
func (s *CarService) GetModels(carMake string) ([]string, error) {
	filters := models.CarFilters{
		Make:  carMake,
		Limit: 1000,
	}

	cars, err := s.repo.GetAll(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}

	// Extract unique models
	modelsMap := make(map[string]bool)
	for _, car := range cars {
		modelsMap[car.Model] = true
	}

	carModels := make([]string, 0, len(modelsMap))
	for model := range modelsMap {
		carModels = append(carModels, model)
	}

	return carModels, nil
}

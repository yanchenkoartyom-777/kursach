package services

import (
	"car-dealership-system/internal/models"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// ImportService handles data import operations
type ImportService struct {
	carService *CarService
	validator  *validator.Validate
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	SuccessCount int      `json:"success_count"`
	ErrorCount   int      `json:"error_count"`
	Errors       []string `json:"errors"`
}

// NewImportService creates a new import service
func NewImportService(carService *CarService) *ImportService {
	return &ImportService{
		carService: carService,
		validator:  validator.New(),
	}
}

// ImportFromJSON imports cars from JSON format
func (s *ImportService) ImportFromJSON(filename string) (*ImportResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var cars []*models.Car
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&cars); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return s.importCars(cars), nil
}

// ImportFromXML imports cars from XML format
func (s *ImportService) ImportFromXML(filename string) (*ImportResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	type CarsXML struct {
		XMLName xml.Name      `xml:"cars"`
		Cars    []*models.Car `xml:"car"`
	}

	var data CarsXML
	decoder := xml.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	return s.importCars(data.Cars), nil
}

// ImportFromCSV imports cars from CSV format
func (s *ImportService) ImportFromCSV(filename string) (*ImportResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Read records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records: %w", err)
	}

	var cars []*models.Car
	for _, record := range records {
		if len(record) < 12 {
			continue
		}

		year, _ := strconv.Atoi(record[3])
		price, _ := strconv.ParseFloat(record[4], 64)
		mileage, _ := strconv.Atoi(record[8])
		isAvailable, _ := strconv.ParseBool(record[11])

		car := &models.Car{
			Make:         record[1],
			Model:        record[2],
			Year:         year,
			Price:        price,
			Color:        record[5],
			EngineType:   record[6],
			Transmission: record[7],
			Mileage:      mileage,
			Description:  record[9],
			ImageURL: func() *string {
				if record[10] != "" {
					s := record[10]
					return &s
				}
				return nil
			}(),
			IsAvailable: isAvailable,
		}

		// Validate car
		if err := s.validator.Struct(car); err != nil {
			// Skip invalid car
			continue
		}

		cars = append(cars, car)
	}

	return s.importCars(cars), nil
}

// importCars imports a slice of cars and returns the result
func (s *ImportService) importCars(cars []*models.Car) *ImportResult {
	result := &ImportResult{
		Errors: []string{},
	}

	for i, car := range cars {
		// Reset ID to create new car
		car.ID = 0

		if err := s.carService.CreateCar(car); err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Car %d (%s %s): %v", i+1, car.Make, car.Model, err))
		} else {
			result.SuccessCount++
		}
	}

	return result
}

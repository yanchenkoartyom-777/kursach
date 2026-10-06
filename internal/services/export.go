package services

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/repositories"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

// ExportService handles data export operations
type ExportService struct {
	carRepo *repositories.CarRepository
}

// NewExportService creates a new export service
func NewExportService(carRepo *repositories.CarRepository) *ExportService {
	return &ExportService{
		carRepo: carRepo,
	}
}

// ExportToJSON exports cars to JSON format
func (s *ExportService) ExportToJSON(filename string, filters models.CarFilters) error {
	cars, err := s.carRepo.GetAll(filters)
	if err != nil {
		return fmt.Errorf("failed to get cars: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(cars); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// ExportToXML exports cars to XML format
func (s *ExportService) ExportToXML(filename string, filters models.CarFilters) error {
	cars, err := s.carRepo.GetAll(filters)
	if err != nil {
		return fmt.Errorf("failed to get cars: %w", err)
	}

	type CarsXML struct {
		XMLName xml.Name      `xml:"cars"`
		Cars    []*models.Car `xml:"car"`
	}

	data := CarsXML{Cars: cars}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	file.WriteString(xml.Header)
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode XML: %w", err)
	}

	return nil
}

// ExportToCSV exports cars to CSV format
func (s *ExportService) ExportToCSV(filename string, filters models.CarFilters) error {
	cars, err := s.carRepo.GetAll(filters)
	if err != nil {
		return fmt.Errorf("failed to get cars: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Make", "Model", "Year", "Price", "Color", "Engine Type", "Transmission", "Mileage", "Description", "Image URL", "Is Available"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, car := range cars {
		record := []string{
			fmt.Sprintf("%d", car.ID),
			car.Make,
			car.Model,
			fmt.Sprintf("%d", car.Year),
			fmt.Sprintf("%.2f", car.Price),
			car.Color,
			car.EngineType,
			car.Transmission,
			fmt.Sprintf("%d", car.Mileage),
			car.Description,
			func() string {
				if car.ImageURL != nil {
					return *car.ImageURL
				}
				return ""
			}(),
			fmt.Sprintf("%t", car.IsAvailable),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

package services

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/repositories"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// ReportService handles report generation
type ReportService struct {
	orderRepo *repositories.OrderRepository
}

// NewReportService creates a new report service
func NewReportService(orderRepo *repositories.OrderRepository) *ReportService {
	return &ReportService{
		orderRepo: orderRepo,
	}
}

// GeneratePDFReport generates a PDF sales report
func (s *ReportService) GeneratePDFReport(filename string, period models.TimePeriod) error {
	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return fmt.Errorf("failed to get statistics: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Sales Report")
	pdf.Ln(12)

	// Period
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 8, fmt.Sprintf("Period: %s to %s",
		period.StartDate.Format("2006-01-02"),
		period.EndDate.Format("2006-01-02")))
	pdf.Ln(10)

	// Summary
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Summary")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(90, 8, "Total Orders:")
	pdf.Cell(100, 8, fmt.Sprintf("%d", stats.TotalOrders))
	pdf.Ln(6)

	pdf.Cell(90, 8, "Total Revenue:")
	pdf.Cell(100, 8, fmt.Sprintf("$%.2f", stats.TotalRevenue))
	pdf.Ln(6)

	pdf.Cell(90, 8, "Average Order:")
	pdf.Cell(100, 8, fmt.Sprintf("$%.2f", stats.AverageOrder))
	pdf.Ln(12)

	// Top Cars
	if len(stats.TopCars) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(190, 8, "Top 10 Cars")
		pdf.Ln(8)

		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(60, 8, "Make & Model")
		pdf.Cell(40, 8, "Orders")
		pdf.Cell(40, 8, "Revenue")
		pdf.Ln(6)

		pdf.SetFont("Arial", "", 10)
		for _, car := range stats.TopCars {
			pdf.Cell(60, 6, fmt.Sprintf("%s %s", car.Make, car.Model))
			pdf.Cell(40, 6, fmt.Sprintf("%d", car.OrderCount))
			pdf.Cell(40, 6, fmt.Sprintf("$%.2f", car.Revenue))
			pdf.Ln(6)
		}
	}

	return pdf.OutputFileAndClose(filename)
}

// GenerateExcelReport generates an Excel sales report
func (s *ReportService) GenerateExcelReport(filename string, period models.TimePeriod) error {
	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return fmt.Errorf("failed to get statistics: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sales Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %w", err)
	}

	// Title
	f.SetCellValue(sheetName, "A1", "Sales Report")
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("Period: %s to %s",
		period.StartDate.Format("2006-01-02"),
		period.EndDate.Format("2006-01-02")))

	// Summary
	f.SetCellValue(sheetName, "A4", "Summary")
	f.SetCellValue(sheetName, "A5", "Total Orders:")
	f.SetCellValue(sheetName, "B5", stats.TotalOrders)
	f.SetCellValue(sheetName, "A6", "Total Revenue:")
	f.SetCellValue(sheetName, "B6", stats.TotalRevenue)
	f.SetCellValue(sheetName, "A7", "Average Order:")
	f.SetCellValue(sheetName, "B7", stats.AverageOrder)

	// Top Cars
	if len(stats.TopCars) > 0 {
		f.SetCellValue(sheetName, "A9", "Top 10 Cars")
		f.SetCellValue(sheetName, "A10", "Make")
		f.SetCellValue(sheetName, "B10", "Model")
		f.SetCellValue(sheetName, "C10", "Orders")
		f.SetCellValue(sheetName, "D10", "Revenue")

		row := 11
		for _, car := range stats.TopCars {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), car.Make)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), car.Model)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), car.OrderCount)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), car.Revenue)
			row++
		}
	}

	// Monthly data
	if len(stats.OrdersByMonth) > 0 {
		f.SetCellValue(sheetName, "A"+fmt.Sprintf("%d", len(stats.TopCars)+13), "Monthly Statistics")
		f.SetCellValue(sheetName, "A"+fmt.Sprintf("%d", len(stats.TopCars)+14), "Month")
		f.SetCellValue(sheetName, "B"+fmt.Sprintf("%d", len(stats.TopCars)+14), "Orders")
		f.SetCellValue(sheetName, "C"+fmt.Sprintf("%d", len(stats.TopCars)+14), "Revenue")

		row := len(stats.TopCars) + 15
		for month, metric := range stats.OrdersByMonth {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), month)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), metric.Count)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), metric.Revenue)
			row++
		}
	}

	f.SetActiveSheet(index)

	return f.SaveAs(filename)
}

// GetReportData returns report data for display
func (s *ReportService) GetReportData(period models.TimePeriod) (*models.OrderStats, error) {
	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return stats, nil
}

// GetDefaultPeriod returns the default reporting period (last 30 days)
func (s *ReportService) GetDefaultPeriod() models.TimePeriod {
	now := time.Now()
	return models.TimePeriod{
		StartDate: now.AddDate(0, 0, -30),
		EndDate:   now,
	}
}

// GetMonthPeriod returns a period for a specific month
func (s *ReportService) GetMonthPeriod(year int, month time.Month) models.TimePeriod {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	return models.TimePeriod{
		StartDate: start,
		EndDate:   end,
	}
}

// GetYearPeriod returns a period for a specific year
func (s *ReportService) GetYearPeriod(year int) models.TimePeriod {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	return models.TimePeriod{
		StartDate: start,
		EndDate:   end,
	}
}

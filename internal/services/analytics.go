package services

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/repositories"
	"fmt"
	"time"
)

// AnalyticsService handles analytics operations
type AnalyticsService struct {
	orderRepo *repositories.OrderRepository
	carRepo   *repositories.CarRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(orderRepo *repositories.OrderRepository, carRepo *repositories.CarRepository) *AnalyticsService {
	return &AnalyticsService{
		orderRepo: orderRepo,
		carRepo:   carRepo,
	}
}

// GetDashboardStats returns statistics for the dashboard
func (s *AnalyticsService) GetDashboardStats() (*models.OrderStats, error) {
	// Get stats for the last 30 days
	period := models.TimePeriod{
		StartDate: time.Now().AddDate(0, 0, -30),
		EndDate:   time.Now(),
	}

	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	return stats, nil
}

// GetMonthlyStats returns statistics grouped by month
func (s *AnalyticsService) GetMonthlyStats(months int) (map[string]models.OrderMetric, error) {
	period := models.TimePeriod{
		StartDate: time.Now().AddDate(0, -months, 0),
		EndDate:   time.Now(),
	}

	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}

	return stats.OrdersByMonth, nil
}

// GetQuarterlyStats returns statistics for the current quarter
func (s *AnalyticsService) GetQuarterlyStats() (*models.OrderStats, error) {
	now := time.Now()
	quarter := (int(now.Month()) - 1) / 3
	startMonth := time.Month(quarter*3 + 1)

	period := models.TimePeriod{
		StartDate: time.Date(now.Year(), startMonth, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   now,
	}

	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get quarterly stats: %w", err)
	}

	return stats, nil
}

// GetYearlyStats returns statistics for the current year
func (s *AnalyticsService) GetYearlyStats() (*models.OrderStats, error) {
	now := time.Now()
	period := models.TimePeriod{
		StartDate: time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   now,
	}

	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly stats: %w", err)
	}

	return stats, nil
}

// GetTopCars returns the top selling cars
func (s *AnalyticsService) GetTopCars(limit int, period models.TimePeriod) ([]models.CarStat, error) {
	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get top cars: %w", err)
	}

	// Limit results
	if limit > 0 && len(stats.TopCars) > limit {
		return stats.TopCars[:limit], nil
	}

	return stats.TopCars, nil
}

// GetRevenueByStatus returns revenue grouped by order status
func (s *AnalyticsService) GetRevenueByStatus(period models.TimePeriod) (map[string]int, error) {
	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue by status: %w", err)
	}

	return stats.OrdersByStatus, nil
}

// GetChartData prepares data for charts
func (s *AnalyticsService) GetChartData(period models.TimePeriod) (map[string]interface{}, error) {
	stats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart data: %w", err)
	}

	// Prepare monthly revenue data for line chart
	monthlyLabels := []string{}
	monthlyRevenue := []float64{}
	monthlyOrders := []int{}

	for month, metric := range stats.OrdersByMonth {
		monthlyLabels = append(monthlyLabels, month)
		monthlyRevenue = append(monthlyRevenue, metric.Revenue)
		monthlyOrders = append(monthlyOrders, metric.Count)
	}

	// Prepare top cars data for bar chart
	topCarLabels := []string{}
	topCarValues := []int{}

	for _, car := range stats.TopCars {
		if len(topCarLabels) >= 10 {
			break
		}
		topCarLabels = append(topCarLabels, fmt.Sprintf("%s %s", car.Make, car.Model))
		topCarValues = append(topCarValues, car.OrderCount)
	}

	// Prepare status data for pie chart
	statusLabels := []string{}
	statusValues := []int{}

	for status, count := range stats.OrdersByStatus {
		statusLabels = append(statusLabels, status)
		statusValues = append(statusValues, count)
	}

	chartData := map[string]interface{}{
		"monthly": map[string]interface{}{
			"labels":  monthlyLabels,
			"revenue": monthlyRevenue,
			"orders":  monthlyOrders,
		},
		"topCars": map[string]interface{}{
			"labels": topCarLabels,
			"values": topCarValues,
		},
		"status": map[string]interface{}{
			"labels": statusLabels,
			"values": statusValues,
		},
	}

	return chartData, nil
}

// ComparePerformance compares performance between two periods
func (s *AnalyticsService) ComparePerformance(period1, period2 models.TimePeriod) (map[string]interface{}, error) {
	stats1, err := s.orderRepo.GetStatistics(period1)
	if err != nil {
		return nil, fmt.Errorf("failed to get period 1 stats: %w", err)
	}

	stats2, err := s.orderRepo.GetStatistics(period2)
	if err != nil {
		return nil, fmt.Errorf("failed to get period 2 stats: %w", err)
	}

	// Calculate changes
	orderChange := float64(0)
	if stats2.TotalOrders > 0 {
		orderChange = ((float64(stats1.TotalOrders) - float64(stats2.TotalOrders)) / float64(stats2.TotalOrders)) * 100
	}

	revenueChange := float64(0)
	if stats2.TotalRevenue > 0 {
		revenueChange = ((stats1.TotalRevenue - stats2.TotalRevenue) / stats2.TotalRevenue) * 100
	}

	comparison := map[string]interface{}{
		"period1": map[string]interface{}{
			"orders":  stats1.TotalOrders,
			"revenue": stats1.TotalRevenue,
			"average": stats1.AverageOrder,
		},
		"period2": map[string]interface{}{
			"orders":  stats2.TotalOrders,
			"revenue": stats2.TotalRevenue,
			"average": stats2.AverageOrder,
		},
		"changes": map[string]interface{}{
			"orders":  orderChange,
			"revenue": revenueChange,
		},
	}

	return comparison, nil
}

// GetEnhancedDashboardStats returns comprehensive statistics for analytics dashboard
func (s *AnalyticsService) GetEnhancedDashboardStats() (*models.EnhancedStats, error) {
	// Get order statistics for last 6 months
	period := models.TimePeriod{
		StartDate: time.Now().AddDate(0, -6, 0),
		EndDate:   time.Now(),
	}

	orderStats, err := s.orderRepo.GetStatistics(period)
	if err != nil {
		return nil, fmt.Errorf("failed to get order stats: %w", err)
	}

	// Get car inventory statistics
	carStats, err := s.carRepo.GetInventoryStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get car inventory stats: %w", err)
	}

	// Get order status breakdown
	statusBreakdown, err := s.orderRepo.GetOrderStatusBreakdown()
	if err != nil {
		return nil, fmt.Errorf("failed to get status breakdown: %w", err)
	}

	// Prepare monthly revenue data
	monthlyRevenue := []models.MonthlyMetric{}
	for month, metric := range orderStats.OrdersByMonth {
		monthlyRevenue = append(monthlyRevenue, models.MonthlyMetric{
			Month:   month,
			Revenue: metric.Revenue,
			Orders:  metric.Count,
		})
	}

	// Get top 5 cars
	topCars := orderStats.TopCars
	if len(topCars) > 5 {
		topCars = topCars[:5]
	}

	enhancedStats := &models.EnhancedStats{
		TotalOrders:    orderStats.TotalOrders,
		TotalRevenue:   orderStats.TotalRevenue,
		AverageOrder:   orderStats.AverageOrder,
		TotalCars:      carStats.TotalCars,
		AvailableCars:  carStats.AvailableCars,
		SoldCars:       carStats.SoldCars,
		OrdersByStatus: statusBreakdown,
		TopCars:        topCars,
		MonthlyRevenue: monthlyRevenue,
	}

	return enhancedStats, nil
}

// GetCarInventoryStats returns car inventory statistics
func (s *AnalyticsService) GetCarInventoryStats() (*models.CarInventoryStats, error) {
	stats, err := s.carRepo.GetInventoryStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get car inventory stats: %w", err)
	}

	return stats, nil
}

// GetOrderStatusBreakdown returns order counts by status
func (s *AnalyticsService) GetOrderStatusBreakdown() (map[string]int, error) {
	breakdown, err := s.orderRepo.GetOrderStatusBreakdown()
	if err != nil {
		return nil, fmt.Errorf("failed to get order status breakdown: %w", err)
	}

	return breakdown, nil
}

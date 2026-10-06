package handlers

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/services"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	carService       *services.CarService
	orderService     *services.OrderService
	analyticsService *services.AnalyticsService
	reportService    *services.ReportService
	exportService    *services.ExportService
	importService    *services.ImportService
}

func NewAdminHandler(
	carService *services.CarService,
	orderService *services.OrderService,
	analyticsService *services.AnalyticsService,
	reportService *services.ReportService,
	exportService *services.ExportService,
	importService *services.ImportService,
) *AdminHandler {
	return &AdminHandler{
		carService:       carService,
		orderService:     orderService,
		analyticsService: analyticsService,
		reportService:    reportService,
		exportService:    exportService,
		importService:    importService,
	}
}

func (h *AdminHandler) ShowAnalytics(c *gin.Context) {
	stats, err := h.analyticsService.GetEnhancedDashboardStats()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin/analytics.html", gin.H{
			"Title": "Аналитика",
			"Error": "Не удалось загрузить статистику",
			"User":  c.MustGet("user"),
		})
		return
	}

	c.HTML(http.StatusOK, "admin/analytics.html", gin.H{
		"Title": "Аналитика",
		"Stats": stats,
		"User":  c.MustGet("user"),
	})
}

func (h *AdminHandler) ShowCars(c *gin.Context) {
	cars, _ := h.carService.GetCars(models.CarFilters{Limit: 100})

	c.HTML(http.StatusOK, "admin/cars.html", gin.H{
		"Title": "Управление автомобилями",
		"Cars":  cars,
		"User":  c.MustGet("user"),
	})
}

func (h *AdminHandler) ShowOrders(c *gin.Context) {
	orders, _ := h.orderService.GetAllOrders(100, 0)

	c.HTML(http.StatusOK, "admin/orders.html", gin.H{
		"Title":  "Управление заказами",
		"Orders": orders,
		"User":   c.MustGet("user"),
	})
}

// UpdateCarImage updates the image URL for a car
func (h *AdminHandler) UpdateCarImage(c *gin.Context) {
	var req struct {
		CarID    uint   `json:"car_id" form:"car_id" binding:"required"`
		ImageURL string `json:"image_url" form:"image_url"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get the car
	car, err := h.carService.GetCarByID(req.CarID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	// Update image URL
	if req.ImageURL != "" {
		car.ImageURL = &req.ImageURL
	} else {
		car.ImageURL = nil
	}

	// Save changes
	if err := h.carService.UpdateCar(car); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update car"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image updated successfully",
	})
}

// ShowCreateCar shows the create car form
func (h *AdminHandler) ShowCreateCar(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/car-form.html", gin.H{
		"Title":  "Добавить автомобиль",
		"User":   c.MustGet("user"),
		"Action": "create",
	})
}

// CreateCar creates a new car
func (h *AdminHandler) CreateCar(c *gin.Context) {
	var car models.Car
	if err := c.ShouldBind(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	car.IsAvailable = true

	if err := h.carService.CreateCar(&car); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Автомобиль успешно добавлен",
		"car_id":  car.ID,
	})
}

// ShowEditCar shows the edit car form
func (h *AdminHandler) ShowEditCar(c *gin.Context) {
	carID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	car, err := h.carService.GetCarByID(uint(carID))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Error": "Автомобиль не найден",
			"User":  c.MustGet("user"),
		})
		return
	}

	c.HTML(http.StatusOK, "admin/car-form.html", gin.H{
		"Title":  "Редактировать автомобиль",
		"User":   c.MustGet("user"),
		"Car":    car,
		"Action": "edit",
	})
}

// UpdateCar updates an existing car
func (h *AdminHandler) UpdateCar(c *gin.Context) {
	carID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var car models.Car
	if err := c.ShouldBind(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	car.ID = uint(carID)

	if err := h.carService.UpdateCar(&car); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Автомобиль успешно обновлён",
	})
}

// DeleteCar deletes a car
func (h *AdminHandler) DeleteCar(c *gin.Context) {
	carID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.carService.DeleteCar(uint(carID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Автомобиль успешно удалён",
	})
}

// ToggleCarAvailability toggles car availability
func (h *AdminHandler) ToggleCarAvailability(c *gin.Context) {
	carID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	car, err := h.carService.GetCarByID(uint(carID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	if car.IsAvailable {
		err = h.carService.MarkAsUnavailable(uint(carID))
	} else {
		err = h.carService.MarkAsAvailable(uint(carID))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"is_available": !car.IsAvailable,
		"message":      "Статус автомобиля изменён",
	})
}

// UpdateOrderStatus updates order status
func (h *AdminHandler) UpdateOrderStatus(c *gin.Context) {
	orderID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := h.orderService.UpdateOrderStatus(uint(orderID), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Статус заказа обновлён",
	})
}

// GetOrderDetails returns order details as JSON
func (h *AdminHandler) GetOrderDetails(c *gin.Context) {
	orderID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	order, err := h.orderService.GetOrderByID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ExportOrdersReport exports orders report
func (h *AdminHandler) ExportOrdersReport(c *gin.Context) {
	format := c.Query("format")

	period := models.TimePeriod{
		StartDate: time.Now().AddDate(0, -6, 0),
		EndDate:   time.Now(),
	}

	filename := fmt.Sprintf("orders_report_%s.%s", time.Now().Format("2006-01-02"), format)

	var err error
	switch format {
	case "pdf":
		err = h.reportService.GeneratePDFReport(filename, period)
	case "excel":
		err = h.reportService.GenerateExcelReport(filename, period)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.File(filename)
	os.Remove(filename)
}

// ShowImportExport shows import/export page
func (h *AdminHandler) ShowImportExport(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/import-export.html", gin.H{
		"Title": "Импорт / Экспорт",
		"User":  c.MustGet("user"),
	})
}

// ExportCars exports cars data
func (h *AdminHandler) ExportCars(c *gin.Context) {
	format := c.Query("format")
	filename := fmt.Sprintf("cars_%s.%s", time.Now().Format("2006-01-02"), format)
	filters := models.CarFilters{Limit: 1000}

	var err error
	switch format {
	case "json":
		err = h.exportService.ExportToJSON(filename, filters)
	case "xml":
		err = h.exportService.ExportToXML(filename, filters)
	case "csv":
		err = h.exportService.ExportToCSV(filename, filters)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.File(filename)
	os.Remove(filename)
}

// ImportCars imports cars data
func (h *AdminHandler) ImportCars(c *gin.Context) {
	format := c.Query("format")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Read file content
	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Save temp file
	tempFile := fmt.Sprintf("temp_%s.%s", time.Now().Format("20060102150405"), format)
	if err := c.SaveUploadedFile(file, tempFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer os.Remove(tempFile)

	var result *services.ImportResult
	switch format {
	case "json":
		result, err = h.importService.ImportFromJSON(tempFile)
	case "xml":
		result, err = h.importService.ImportFromXML(tempFile)
	case "csv":
		result, err = h.importService.ImportFromCSV(tempFile)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Данные успешно импортированы",
		"imported": result.SuccessCount,
	})
}

// ExportCarsReport exports cars report
func (h *AdminHandler) ExportCarsReport(c *gin.Context) {
	format := c.Query("format")
	filename := fmt.Sprintf("cars_report_%s.%s", time.Now().Format("2006-01-02"), format)

	period := models.TimePeriod{
		StartDate: time.Now().AddDate(-1, 0, 0),
		EndDate:   time.Now(),
	}

	var err error
	switch format {
	case "pdf":
		err = h.reportService.GeneratePDFReport(filename, period)
	case "excel":
		err = h.reportService.GenerateExcelReport(filename, period)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.File(filename)
	os.Remove(filename)
}

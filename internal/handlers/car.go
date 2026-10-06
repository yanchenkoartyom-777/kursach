package handlers

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	carService *services.CarService
}

func NewCarHandler(carService *services.CarService) *CarHandler {
	return &CarHandler{carService: carService}
}

func (h *CarHandler) ShowCatalog(c *gin.Context) {
	filters := models.CarFilters{
		Limit: 50,
	}

	// Get search and sort parameters
	search := c.Query("search")
	sort := c.Query("sort")

	// Parse sort parameter
	switch sort {
	case "price_asc":
		filters.SortBy = "price"
		filters.SortOrder = "asc"
	case "price_desc":
		filters.SortBy = "price"
		filters.SortOrder = "desc"
	case "year_desc":
		filters.SortBy = "year"
		filters.SortOrder = "desc"
	case "year_asc":
		filters.SortBy = "year"
		filters.SortOrder = "asc"
	default:
		filters.SortBy = "created_at"
		filters.SortOrder = "desc"
	}

	// Get user if authenticated (optional)
	user, _ := c.Get("user")

	var cars []*models.Car
	var err error

	// Check if user is admin - admins see all cars, clients see only available
	if user != nil {
		if userMap, ok := user.(map[string]interface{}); ok {
			if role, ok := userMap["Role"].(string); ok && role == "admin" {
				// Admin sees all cars - don't filter by availability
				filters.IsAvailable = nil
			}
		}
	}

	if search != "" {
		cars, err = h.carService.SearchCars(search, filters)
	} else {
		if filters.IsAvailable == nil {
			// Admin - get all cars
			cars, err = h.carService.GetCars(filters)
		} else {
			// Client - get only available cars
			cars, err = h.carService.GetAvailableCars(filters)
		}
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Ошибка загрузки автомобилей",
			"User":  user,
		})
		return
	}

	c.HTML(http.StatusOK, "cars/catalog.html", gin.H{
		"Title":  "Каталог",
		"Cars":   cars,
		"User":   user,
		"Search": search,
		"Sort":   sort,
	})
}

func (h *CarHandler) ShowDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	car, err := h.carService.GetCarByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Error": "Автомобиль не найден",
		})
		return
	}

	// Get user if authenticated (optional)
	user, _ := c.Get("user")

	c.HTML(http.StatusOK, "cars/detail.html", gin.H{
		"Title": car.Make + " " + car.Model,
		"Car":   car,
		"User":  user,
	})
}

package handlers

import (
	"car-dealership-system/internal/middleware"
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) ShowCart(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	cart, _ := h.orderService.GetCart(userID)

	c.HTML(http.StatusOK, "orders/cart.html", gin.H{
		"Title": "Корзина",
		"Cart":  cart,
		"User":  c.MustGet("user"),
	})
}

func (h *OrderHandler) AddToCart(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	carID, _ := strconv.ParseUint(c.PostForm("car_id"), 10, 32)

	if err := h.orderService.AddToCart(userID, uint(carID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *OrderHandler) RemoveFromCart(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	carID, _ := strconv.ParseUint(c.PostForm("car_id"), 10, 32)

	h.orderService.RemoveFromCart(userID, uint(carID))
	c.Redirect(http.StatusFound, "/cart")
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	order, err := h.orderService.CreateOrder(userID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "orders/cart.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "orders/checkout.html", gin.H{
		"Title": "Заказ оформлен",
		"Order": order,
		"User":  c.MustGet("user"),
	})
}

func (h *OrderHandler) ShowHistory(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	orders, _ := h.orderService.GetOrderHistory(userID)

	c.HTML(http.StatusOK, "orders/history.html", gin.H{
		"Title":  "История заказов",
		"Orders": orders,
		"User":   c.MustGet("user"),
	})
}

func (h *OrderHandler) GetCartCount(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	count, _ := h.orderService.GetCartItemCount(userID)

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// ShowMyOrders shows customer's own orders or all orders for admin
func (h *OrderHandler) ShowMyOrders(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	userRole, _ := middleware.GetUserRole(c)

	var orders []*models.Order
	var err error

	// If user is admin, show all orders
	if userRole == "admin" {
		orders, err = h.orderService.GetAllOrders(100, 0)
	} else {
		orders, err = h.orderService.GetOrderHistory(userID)
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError, "orders/my-orders.html", gin.H{
			"Title":  "Мои заказы",
			"Error":  "Не удалось загрузить заказы",
			"User":   c.MustGet("user"),
			"Orders": []*models.Order{},
		})
		return
	}

	c.HTML(http.StatusOK, "orders/my-orders.html", gin.H{
		"Title":  "Мои заказы",
		"Orders": orders,
		"User":   c.MustGet("user"),
	})
}

package main

import (
	"fmt"
	"log"
	"os"

	"car-dealership-system/internal/config"
	"car-dealership-system/internal/handlers"
	"car-dealership-system/internal/middleware"
	"car-dealership-system/internal/repositories"
	"car-dealership-system/internal/services"
	"car-dealership-system/pkg/database"
	"car-dealership-system/pkg/renderer"
	"car-dealership-system/pkg/session"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	configPath := "config.yaml"
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		configPath = path
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	logger := setupLogger(cfg.Logging)
	logger.Info("Starting Car Dealership System")

	// Connect to database
	db, err := database.Connect(cfg.Database.GetDSN())
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("Connected to database successfully")

	// Note: Migrations are now managed via golang-migrate CLI
	// Run: make migrate-up to apply migrations
	// See MIGRATIONS.md for more information

	// Initialize session manager
	sessionManager, err := session.NewManager(cfg.Session.Secret, cfg.Session.MaxAge)
	if err != nil {
		logger.Fatalf("Failed to create session manager: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	twoFactorRepo := repositories.NewTwoFactorRepository(db)
	passwordResetRepo := repositories.NewPasswordResetRepository(db)
	carRepo := repositories.NewCarRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Initialize services
	emailService := services.NewEmailService(cfg.Email)
	codeExpiration, _ := cfg.GetCodeExpirationDuration()
	twoFactorService := services.NewTwoFactorService(twoFactorRepo, cfg.TwoFactor.CodeLength, codeExpiration)
	authService := services.NewAuthService(userRepo, passwordResetRepo, twoFactorService, emailService, cfg.TwoFactor.Enabled)
	carService := services.NewCarService(carRepo)
	orderService := services.NewOrderService(orderRepo, carRepo, emailService, userRepo)
	exportService := services.NewExportService(carRepo)
	importService := services.NewImportService(carService)
	reportService := services.NewReportService(orderRepo)
	analyticsService := services.NewAnalyticsService(orderRepo, carRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, sessionManager)
	carHandler := handlers.NewCarHandler(carService)
	orderHandler := handlers.NewOrderHandler(orderService)
	adminHandler := handlers.NewAdminHandler(carService, orderService, analyticsService, reportService, exportService, importService)

	// Setup Gin
	if cfg.Logging.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.NoCache())                    // Prevent browser caching
	router.Use(middleware.OptionalAuth(sessionManager)) // Try to authenticate but don't require

	// Use custom template renderer that loads templates on-demand
	router.HTMLRender = &renderer.TemplateRenderer{
		TemplatesPath: cfg.Server.TemplatesPath,
	}
	logger.Info("Template renderer configured")

	// Serve static files
	router.Static("/static", cfg.Server.StaticPath)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Favicon handler (prevent 404 errors)
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // No content
	})

	// Public routes
	router.GET("/", func(c *gin.Context) {
		user, _ := c.Get("user")
		c.HTML(200, "home/index.html", gin.H{
			"Title": "Главная",
			"User":  user,
		})
	})
	router.GET("/login", authHandler.ShowLogin)
	router.POST("/login", authHandler.Login)
	router.GET("/register", authHandler.ShowRegister)
	router.POST("/register", authHandler.Register)
	router.GET("/verify", authHandler.ShowVerify)
	router.POST("/verify", authHandler.Verify)
	router.GET("/forgot-password", authHandler.ShowForgotPassword)
	router.POST("/forgot-password", authHandler.ForgotPassword)
	router.GET("/reset-password", authHandler.ShowResetPassword)
	router.POST("/reset-password", authHandler.ResetPassword)
	router.GET("/logout", authHandler.Logout)
	router.GET("/cars", carHandler.ShowCatalog)
	router.GET("/cars/:id", carHandler.ShowDetail)

	// Protected routes (require authentication)
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(sessionManager))
	{
		auth.GET("/profile", authHandler.ShowProfile)
		auth.GET("/cart", orderHandler.ShowCart)
		auth.POST("/cart/add", orderHandler.AddToCart)
		auth.POST("/cart/remove", orderHandler.RemoveFromCart)
		auth.GET("/cart/count", orderHandler.GetCartCount)
		auth.POST("/orders", orderHandler.CreateOrder)
		auth.GET("/orders", orderHandler.ShowMyOrders)
		auth.GET("/orders/history", orderHandler.ShowHistory)
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(middleware.RequireAdmin(sessionManager))
	{
		admin.GET("/analytics", adminHandler.ShowAnalytics)
		admin.GET("/dashboard", func(c *gin.Context) {
			c.Redirect(301, "/admin/analytics")
		})
		admin.GET("/cars", adminHandler.ShowCars)
		admin.GET("/cars/create", adminHandler.ShowCreateCar)
		admin.POST("/cars/create", adminHandler.CreateCar)
		admin.GET("/cars/:id/edit", adminHandler.ShowEditCar)
		admin.POST("/cars/:id/update", adminHandler.UpdateCar)
		admin.POST("/cars/:id/delete", adminHandler.DeleteCar)
		admin.POST("/cars/:id/toggle", adminHandler.ToggleCarAvailability)
		admin.POST("/cars/update-image", adminHandler.UpdateCarImage)
		admin.GET("/orders", adminHandler.ShowOrders)
		admin.POST("/orders/:id/status", adminHandler.UpdateOrderStatus)
		admin.GET("/orders/:id/details", adminHandler.GetOrderDetails)
		admin.GET("/reports/orders", adminHandler.ExportOrdersReport)
		admin.GET("/import-export", adminHandler.ShowImportExport)
		admin.GET("/export/cars", adminHandler.ExportCars)
		admin.POST("/import/cars", adminHandler.ImportCars)
		admin.GET("/reports/cars", adminHandler.ExportCarsReport)
	}

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Server starting on %s", addr)

	if err := router.Run(addr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

func setupLogger(cfg config.LoggingConfig) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set log format
	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return logger
}

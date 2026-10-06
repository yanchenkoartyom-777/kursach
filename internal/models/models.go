package models

import "time"

// User represents a user in the system
type User struct {
	ID           uint      `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FirstName    string    `json:"first_name" db:"first_name" validate:"required"`
	LastName     string    `json:"last_name" db:"last_name" validate:"required"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Car represents a car in the dealership
type Car struct {
	ID           uint      `json:"id" db:"id"`
	Make         string    `json:"make" db:"make" validate:"required"`
	Model        string    `json:"model" db:"model" validate:"required"`
	Year         int       `json:"year" db:"year" validate:"required,min=1900,max=2100"`
	Price        float64   `json:"price" db:"price" validate:"required,min=0"`
	Color        string    `json:"color" db:"color"`
	EngineType   string    `json:"engine_type" db:"engine_type"`
	Transmission string    `json:"transmission" db:"transmission"`
	Mileage      int       `json:"mileage" db:"mileage" validate:"min=0"`
	Description  string    `json:"description" db:"description"`
	ImageURL     *string   `json:"image_url,omitempty" db:"image_url"`
	IsAvailable  bool      `json:"is_available" db:"is_available"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Order represents a customer order
type Order struct {
	ID          uint        `json:"id" db:"id"`
	UserID      uint        `json:"user_id" db:"user_id"`
	User        *User       `json:"user,omitempty"`
	TotalAmount float64     `json:"total_amount" db:"total_amount"`
	Status      string      `json:"status" db:"status"`
	Items       []OrderItem `json:"items,omitempty"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uint      `json:"id" db:"id"`
	OrderID   uint      `json:"order_id" db:"order_id"`
	CarID     uint      `json:"car_id" db:"car_id"`
	Car       *Car      `json:"car,omitempty"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        uint      `json:"id" db:"id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	CarID     uint      `json:"car_id" db:"car_id"`
	Car       *Car      `json:"car,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TwoFactorCode represents a two-factor authentication code
type TwoFactorCode struct {
	ID        uint      `json:"id" db:"id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	Code      string    `json:"code" db:"code"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	IsUsed    bool      `json:"is_used" db:"is_used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        uint      `json:"id" db:"id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	IsUsed    bool      `json:"is_used" db:"is_used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" form:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" form:"first_name" validate:"required"`
	LastName  string `json:"last_name" form:"last_name" validate:"required"`
}

// VerifyRequest represents a 2FA verification request
type VerifyRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
	Code  string `json:"code" form:"code" validate:"required,len=6"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

// ResetPasswordRequest represents a reset password request
type ResetPasswordRequest struct {
	Token    string `json:"token" form:"token" validate:"required"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

// LoginResponse represents the response after initial login
type LoginResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
	User    *User  `json:"user,omitempty"`
	Skip2FA bool   `json:"skip_2fa,omitempty"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	User *User `json:"user"`
}

// CarFilters represents filters for car search
type CarFilters struct {
	Make         string
	Model        string
	MinYear      int
	MaxYear      int
	MinPrice     float64
	MaxPrice     float64
	EngineType   string
	Transmission string
	IsAvailable  *bool
	SortBy       string // price, year, created_at
	SortOrder    string // asc, desc
	Limit        int
	Offset       int
}

// SearchQuery represents a search query
type SearchQuery struct {
	Query   string
	Filters CarFilters
}

// Cart represents a shopping cart with items
type Cart struct {
	Items       []CartItem `json:"items"`
	TotalAmount float64    `json:"total_amount"`
	ItemCount   int        `json:"item_count"`
}

// OrderStats represents order statistics
type OrderStats struct {
	TotalOrders    int                    `json:"total_orders"`
	TotalRevenue   float64                `json:"total_revenue"`
	AverageOrder   float64                `json:"average_order"`
	OrdersByStatus map[string]int         `json:"orders_by_status"`
	OrdersByMonth  map[string]OrderMetric `json:"orders_by_month"`
	TopCars        []CarStat              `json:"top_cars"`
}

// OrderMetric represents metrics for a time period
type OrderMetric struct {
	Count   int     `json:"count"`
	Revenue float64 `json:"revenue"`
}

// CarStat represents statistics for a car
type CarStat struct {
	CarID      uint    `json:"car_id"`
	Make       string  `json:"make"`
	Model      string  `json:"model"`
	OrderCount int     `json:"order_count"`
	Revenue    float64 `json:"revenue"`
}

// TimePeriod represents a time period for statistics
type TimePeriod struct {
	StartDate time.Time
	EndDate   time.Time
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// CarInventoryStats represents car inventory statistics
type CarInventoryStats struct {
	TotalCars       int `json:"total_cars"`
	AvailableCars   int `json:"available_cars"`
	SoldCars        int `json:"sold_cars"`
	UnavailableCars int `json:"unavailable_cars"`
}

// EnhancedStats represents comprehensive statistics for analytics dashboard
type EnhancedStats struct {
	// Order Statistics
	TotalOrders  int     `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
	AverageOrder float64 `json:"average_order"`

	// Car Statistics
	TotalCars     int `json:"total_cars"`
	AvailableCars int `json:"available_cars"`
	SoldCars      int `json:"sold_cars"`

	// Status Breakdown
	OrdersByStatus map[string]int `json:"orders_by_status"`

	// Top Cars
	TopCars []CarStat `json:"top_cars"`

	// Monthly Data
	MonthlyRevenue []MonthlyMetric `json:"monthly_revenue"`
}

// MonthlyMetric represents metrics for a specific month
type MonthlyMetric struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
	Orders  int     `json:"orders"`
}

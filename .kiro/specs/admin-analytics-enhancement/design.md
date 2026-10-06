# Design Document

## Overview

This design enhances the admin panel by transforming the Dashboard into a comprehensive Analytics page with detailed statistics and visualizations. It also restructures the order management to show all orders in the admin panel while maintaining a separate "My Orders" view for regular customers.

## Architecture

### Component Structure

```
Admin Panel
├── Analytics Dashboard (renamed from Dashboard)
│   ├── Statistics Cards
│   ├── Revenue Trend Chart
│   ├── Top Cars Chart
│   └── Order Status Chart
├── Order Management
│   ├── All Orders View (Admin)
│   └── Order Filtering & Sorting
└── Client Interface
    └── My Orders View (Customer)
```

### Route Changes

**Before:**
- `/admin/dashboard` - Dashboard
- `/admin/analytics` - Analytics
- `/admin/orders` - Orders

**After:**
- `/admin/analytics` - Main Analytics Dashboard (replaces dashboard)
- `/admin/orders` - All Orders Management
- `/orders` - Customer's Own Orders

## Components and Interfaces

### 1. Analytics Dashboard Component

**Location:** `templates/admin/analytics.html`

**Statistics Cards:**
- Total Orders (count)
- Total Revenue (currency formatted)
- Average Order Value (currency formatted)
- Total Cars in Inventory
- Available Cars
- Sold Cars
- Order Status Breakdown (pending, processing, completed, cancelled)

**Charts:**
1. **Revenue Trend Line Chart** - Shows revenue over last 6 months
2. **Top 5 Cars Bar Chart** - Best-selling car models
3. **Order Status Pie Chart** - Distribution of orders by status

### 2. Enhanced Analytics Service

**Location:** `internal/services/analytics.go`

**New Methods:**
```go
// GetEnhancedDashboardStats returns comprehensive statistics
func (s *AnalyticsService) GetEnhancedDashboardStats() (*models.EnhancedStats, error)

// GetCarInventoryStats returns car inventory statistics
func (s *AnalyticsService) GetCarInventoryStats() (*models.CarInventoryStats, error)

// GetOrderStatusBreakdown returns order counts by status
func (s *AnalyticsService) GetOrderStatusBreakdown() (map[string]int, error)
```

### 3. Navigation Menu Updates

**Location:** `templates/layouts/admin.html`

**Changes:**
- Remove "Dashboard" menu item
- Rename first item to "Аналитика" with route `/admin/analytics`
- Keep "Заказы" for order management
- Remove "Отчеты" as it's integrated into Analytics
- Update icons and ordering

### 4. Order Management Enhancement

**Admin View (`/admin/orders`):**
- Shows ALL orders from ALL customers
- Displays customer information (name, email)
- Allows status updates
- Provides filtering by status
- Enables sorting by date, amount, status

**Customer View (`/orders`):**
- Shows only logged-in customer's orders
- Displays order details and status
- Allows viewing order items

## Data Models

### Enhanced Statistics Model

```go
type EnhancedStats struct {
    // Order Statistics
    TotalOrders   int
    TotalRevenue  float64
    AverageOrder  float64
    
    // Car Statistics
    TotalCars     int
    AvailableCars int
    SoldCars      int
    
    // Status Breakdown
    OrdersByStatus map[string]int
    
    // Top Cars
    TopCars []CarStat
    
    // Monthly Data
    MonthlyRevenue []MonthlyMetric
}

type CarInventoryStats struct {
    TotalCars     int
    AvailableCars int
    SoldCars      int
    UnavailableCars int
}

type MonthlyMetric struct {
    Month    string
    Revenue  float64
    Orders   int
}
```

## Implementation Details

### 1. Analytics Handler Updates

**File:** `internal/handlers/admin.go`

**Changes:**
- Rename `ShowDashboard` to `ShowAnalytics`
- Update route from `/admin/dashboard` to `/admin/analytics`
- Enhance data fetching to include car inventory stats
- Add order status breakdown

### 2. Repository Enhancements

**File:** `internal/repositories/car.go`

**New Methods:**
```go
// GetInventoryStats returns car inventory statistics
func (r *CarRepository) GetInventoryStats() (*models.CarInventoryStats, error)

// GetSoldCarsCount returns count of sold cars
func (r *CarRepository) GetSoldCarsCount() (int, error)
```

**File:** `internal/repositories/order.go`

**New Methods:**
```go
// GetOrderStatusBreakdown returns order counts grouped by status
func (r *OrderRepository) GetOrderStatusBreakdown() (map[string]int, error)
```

### 3. Template Structure

**Analytics Dashboard Template:**
```html
<!-- Statistics Cards Row -->
<div class="row mb-4">
    <!-- 6 stat cards: Orders, Revenue, Avg Order, Total Cars, Available, Sold -->
</div>

<!-- Order Status Breakdown Row -->
<div class="row mb-4">
    <!-- 4 cards: Pending, Processing, Completed, Cancelled -->
</div>

<!-- Charts Row 1 -->
<div class="row mb-4">
    <!-- Revenue Trend Chart (full width) -->
</div>

<!-- Charts Row 2 -->
<div class="row">
    <!-- Top 5 Cars Chart (half width) -->
    <!-- Order Status Pie Chart (half width) -->
</div>
```

### 4. Chart Configuration

**Chart.js Implementation:**

1. **Revenue Trend Chart** (Line)
   - X-axis: Last 6 months
   - Y-axis: Revenue in dollars
   - Data: Monthly revenue totals

2. **Top Cars Chart** (Bar)
   - X-axis: Car models
   - Y-axis: Number of sales
   - Data: Top 5 best-selling cars

3. **Order Status Chart** (Pie)
   - Segments: Order statuses
   - Data: Count per status
   - Colors: Status-specific (pending=yellow, processing=blue, completed=green, cancelled=red)

## Error Handling

### Analytics Service Errors
- Return empty stats with zero values if database query fails
- Log errors for monitoring
- Display user-friendly message in UI

### Order Management Errors
- Handle missing orders gracefully
- Validate status updates before applying
- Show error notifications for failed operations

## Testing Strategy

### Unit Tests
- Test analytics service methods for correct calculations
- Test repository methods for data retrieval
- Test status breakdown logic

### Integration Tests
- Test analytics dashboard rendering with real data
- Test order filtering and sorting
- Test navigation menu updates

### Manual Testing
- Verify all statistics display correctly
- Verify charts render with accurate data
- Verify order management shows all orders for admin
- Verify customer sees only their orders
- Verify navigation menu changes

## Migration Path

1. Update navigation menu template
2. Rename dashboard route to analytics
3. Enhance analytics service with new methods
4. Update analytics template with new statistics
5. Add chart implementations
6. Test all functionality
7. Deploy changes

## Performance Considerations

- Cache analytics statistics for 5 minutes to reduce database load
- Use database indexes on order status and created_at fields
- Limit chart data to reasonable time periods (6 months default)
- Paginate order lists for large datasets

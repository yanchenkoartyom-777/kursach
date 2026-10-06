# Implementation Plan

- [x] 1. Update navigation menu and routing


  - Update admin layout template to rename "Dashboard" to "Аналитика"
  - Change route from `/admin/dashboard` to `/admin/analytics`
  - Remove "Отчеты" menu item
  - Update menu item ordering and icons
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 5.1, 5.2, 5.3, 5.4, 5.5, 5.6_

- [x] 2. Enhance analytics service with comprehensive statistics


  - [x] 2.1 Add car inventory statistics methods to car repository


    - Implement `GetInventoryStats()` method to return total, available, and sold car counts
    - Implement `GetSoldCarsCount()` method
    - _Requirements: 2.4, 2.5, 2.6_
  
  - [x] 2.2 Add order status breakdown method to order repository


    - Implement `GetOrderStatusBreakdown()` to return counts grouped by status
    - _Requirements: 2.7_
  
  - [x] 2.3 Create enhanced analytics service methods


    - Implement `GetEnhancedDashboardStats()` combining order and car statistics
    - Implement `GetCarInventoryStats()` for car-specific metrics
    - Implement `GetOrderStatusBreakdown()` wrapper method
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_

- [x] 3. Update admin handler for analytics dashboard


  - Rename `ShowDashboard` method to `ShowAnalytics`
  - Update route handler from `/admin/dashboard` to `/admin/analytics`
  - Fetch enhanced statistics including car inventory and status breakdown
  - Pass comprehensive data to template
  - _Requirements: 1.2, 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_

- [x] 4. Create enhanced analytics dashboard template


  - [x] 4.1 Add statistics cards section


    - Create 6 stat cards: Total Orders, Revenue, Average Order, Total Cars, Available Cars, Sold Cars
    - Format currency values properly
    - Add appropriate icons and styling
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_
  
  - [x] 4.2 Add order status breakdown cards

    - Create 4 status cards: Pending, Processing, Completed, Cancelled
    - Display counts for each status
    - Add status-specific colors
    - _Requirements: 2.7_
  
  - [x] 4.3 Implement revenue trend chart

    - Add Chart.js line chart for monthly revenue
    - Configure chart with last 6 months data
    - Add proper labels and formatting
    - _Requirements: 2.9_
  
  - [x] 4.4 Implement top cars bar chart

    - Add Chart.js bar chart for top 5 cars
    - Display car make/model and sales count
    - Configure horizontal bar layout
    - _Requirements: 2.8_
  
  - [x] 4.5 Implement order status pie chart

    - Add Chart.js pie chart for order status distribution
    - Use status-specific colors
    - Add percentage labels
    - _Requirements: 2.10_

- [x] 5. Update order management for admin view


  - Ensure admin orders page shows all orders from all customers
  - Verify customer name and email display in orders table
  - Confirm filtering and sorting functionality works
  - Update page title to "Заказы" (remove "Управление")
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

- [x] 6. Create customer orders view


  - [x] 6.1 Create customer orders handler


    - Implement handler method to show only logged-in customer's orders
    - Add route `/orders` for customer access
    - Filter orders by current user ID
    - _Requirements: 4.1, 4.2_
  
  - [x] 6.2 Create customer orders template


    - Create template showing customer's own orders
    - Display order details, items, total amount, and status
    - Add order tracking information
    - Style consistently with client interface
    - _Requirements: 4.2, 4.3_
  
  - [x] 6.3 Add navigation link for customer orders


    - Add "Мои заказы" link to client navigation menu
    - Ensure link is visible only to authenticated customers
    - _Requirements: 4.4_

- [x] 7. Update routing configuration



  - Update main router to map `/admin/analytics` instead of `/admin/dashboard`
  - Add `/orders` route for customer orders view
  - Ensure proper middleware for authentication and authorization
  - Remove or redirect old `/admin/dashboard` route
  - _Requirements: 1.2, 4.1_

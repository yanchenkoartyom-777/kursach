# Requirements Document

## Introduction

This feature enhances the admin panel by renaming "Dashboard" to "Analytics" and providing comprehensive analytics with detailed statistics. Additionally, it removes "My Orders" from the client interface and replaces it with an "All Orders" view in the admin panel where administrators can see all orders from all users.

## Glossary

- **Admin Panel**: The administrative interface accessible only to users with admin role
- **Analytics Dashboard**: The main analytics page showing comprehensive statistics and charts
- **Order Management System**: The system component that handles viewing and managing all customer orders
- **Client Interface**: The user-facing interface for regular customers
- **Navigation Menu**: The sidebar menu in the admin panel

## Requirements

### Requirement 1

**User Story:** As an administrator, I want the main dashboard to be renamed to "Analytics" so that it better reflects its purpose of providing business insights

#### Acceptance Criteria

1. WHEN THE Admin Panel loads, THE Navigation Menu SHALL display "Аналитика" instead of "Dashboard"
2. THE Analytics Dashboard SHALL be accessible at the route "/admin/analytics"
3. THE Analytics Dashboard page title SHALL display "Аналитика" in Russian
4. THE Navigation Menu icon for Analytics SHALL use "bi-graph-up" Bootstrap icon

### Requirement 2

**User Story:** As an administrator, I want to see comprehensive analytics with detailed statistics so that I can make informed business decisions

#### Acceptance Criteria

1. THE Analytics Dashboard SHALL display total orders count with visual card presentation
2. THE Analytics Dashboard SHALL display total revenue with currency formatting
3. THE Analytics Dashboard SHALL display average order value with currency formatting
4. THE Analytics Dashboard SHALL display total number of cars in inventory
5. THE Analytics Dashboard SHALL display number of available cars
6. THE Analytics Dashboard SHALL display number of sold cars
7. THE Analytics Dashboard SHALL display order status breakdown showing counts for pending, processing, completed, and cancelled orders
8. THE Analytics Dashboard SHALL display top 5 best-selling car models with sales count
9. THE Analytics Dashboard SHALL display revenue trend chart for the last 6 months
10. THE Analytics Dashboard SHALL display orders by status pie chart

### Requirement 3

**User Story:** As an administrator, I want to view all orders from all customers in one place so that I can manage the entire order pipeline

#### Acceptance Criteria

1. THE Order Management System SHALL display all orders from all users in the admin orders page
2. WHEN THE administrator accesses "/admin/orders", THE Order Management System SHALL show orders from all customers
3. THE Order Management System SHALL display customer name and email for each order
4. THE Order Management System SHALL allow filtering orders by status
5. THE Order Management System SHALL allow sorting orders by date, amount, or status
6. THE Order Management System SHALL display order statistics cards showing counts by status

### Requirement 4

**User Story:** As a regular customer, I want to see my orders in a dedicated section so that I can track my purchases

#### Acceptance Criteria

1. THE Client Interface SHALL provide a "Мои заказы" page accessible at "/orders"
2. WHEN THE customer accesses "/orders", THE Client Interface SHALL display only orders belonging to that customer
3. THE Client Interface SHALL display order details including items, total amount, and status
4. THE Navigation Menu in client interface SHALL include link to "Мои заказы"

### Requirement 5

**User Story:** As an administrator, I want the navigation menu to clearly distinguish between different admin functions so that I can quickly access the tools I need

#### Acceptance Criteria

1. THE Navigation Menu SHALL display "Аналитика" as the first menu item after the admin panel header
2. THE Navigation Menu SHALL display "Автомобили" for car management
3. THE Navigation Menu SHALL display "Заказы" for order management
4. THE Navigation Menu SHALL display "Импорт/Экспорт" for data operations
5. THE Navigation Menu SHALL remove the separate "Dashboard" entry
6. THE Navigation Menu SHALL remove the "Отчеты" entry as reports are integrated into Analytics

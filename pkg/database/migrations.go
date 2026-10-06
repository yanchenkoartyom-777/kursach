package database

import (
	"database/sql"
	"fmt"
)

// RunMigrations applies all database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsersTable,
		createTwoFactorCodesTable,
		createCarsTable,
		createOrdersTable,
		createOrderItemsTable,
		createCartItemsTable,
		createIndexes,
		seedData,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) DEFAULT 'client',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createTwoFactorCodesTable = `
CREATE TABLE IF NOT EXISTS two_factor_codes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createCarsTable = `
CREATE TABLE IF NOT EXISTS cars (
    id SERIAL PRIMARY KEY,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    color VARCHAR(50),
    engine_type VARCHAR(50),
    transmission VARCHAR(50),
    mileage INTEGER DEFAULT 0,
    description TEXT,
    image_url VARCHAR(500),
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createOrdersTable = `
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createOrderItemsTable = `
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    car_id INTEGER REFERENCES cars(id),
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createCartItemsTable = `
CREATE TABLE IF NOT EXISTS cart_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    car_id INTEGER REFERENCES cars(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, car_id)
);
`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_cars_make_model ON cars(make, model);
CREATE INDEX IF NOT EXISTS idx_cars_price ON cars(price);
CREATE INDEX IF NOT EXISTS idx_cars_year ON cars(year);
CREATE INDEX IF NOT EXISTS idx_cars_available ON cars(is_available);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_two_factor_user_id ON two_factor_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_two_factor_expires_at ON two_factor_codes(expires_at);
CREATE INDEX IF NOT EXISTS idx_cart_user_id ON cart_items(user_id);
`

const seedData = `
-- Seed data for development (Note: These password hashes are examples and should be regenerated)
-- Insert sample admin user (password will be set during first run)
INSERT INTO users (email, password_hash, first_name, last_name, role)
VALUES ('admin@cardealership.com', '$2a$10$placeholder', 'Admin', 'User', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Insert sample cars
INSERT INTO cars (make, model, year, price, color, engine_type, transmission, mileage, description, is_available)
VALUES 
    ('Toyota', 'Camry', 2023, 28500.00, 'Silver', 'Hybrid', 'Automatic', 5000, 'Надежный седан с гибридным двигателем', true),
    ('Honda', 'Accord', 2023, 27800.00, 'Black', 'Gasoline', 'Automatic', 3000, 'Комфортный семейный автомобиль', true),
    ('BMW', '3 Series', 2022, 45000.00, 'White', 'Gasoline', 'Automatic', 15000, 'Премиальный спортивный седан', true),
    ('Mercedes-Benz', 'C-Class', 2022, 48500.00, 'Blue', 'Diesel', 'Automatic', 12000, 'Элегантный и мощный автомобиль', true),
    ('Audi', 'A4', 2023, 42000.00, 'Gray', 'Gasoline', 'Automatic', 8000, 'Современный дизайн и технологии', true),
    ('Tesla', 'Model 3', 2023, 52000.00, 'Red', 'Electric', 'Automatic', 2000, 'Электромобиль будущего', true),
    ('Ford', 'Mustang', 2022, 38000.00, 'Yellow', 'Gasoline', 'Manual', 10000, 'Легендарный американский маслкар', true),
    ('Volkswagen', 'Passat', 2023, 32000.00, 'Silver', 'Diesel', 'Automatic', 6000, 'Практичный и экономичный седан', true),
    ('Mazda', 'CX-5', 2023, 35000.00, 'Blue', 'Gasoline', 'Automatic', 4000, 'Стильный кроссовер для семьи', true),
    ('Hyundai', 'Sonata', 2023, 26500.00, 'White', 'Gasoline', 'Automatic', 7000, 'Доступный и надежный седан', true)
ON CONFLICT DO NOTHING;
`

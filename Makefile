# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=car_dealership
DB_USER=postgres
DB_PASSWORD=postgres
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Migration configuration
MIGRATIONS_PATH=./migrations

.PHONY: help migrate-up migrate-down migrate-reset migrate-force migrate-version migrate-create db-status

help: ## Show this help message
	@echo "Available commands:"
	@echo "  make migrate-up       - Apply all pending migrations"
	@echo "  make migrate-down     - Rollback last migration"
	@echo "  make migrate-reset    - Rollback all migrations and reapply"
	@echo "  make migrate-force V=N - Force set migration version to N"
	@echo "  make migrate-version  - Show current migration version"
	@echo "  make migrate-create NAME=name - Create new migration files"
	@echo "  make db-status        - Show database status"
	@echo ""
	@echo "Examples:"
	@echo "  make migrate-up"
	@echo "  make migrate-down"
	@echo "  make migrate-reset"
	@echo "  make migrate-create NAME=add_user_avatar"

migrate-up: ## Apply all pending migrations
	@echo "Applying migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up
	@echo "✓ Migrations applied successfully"

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1
	@echo "✓ Migration rolled back successfully"

migrate-reset: ## Rollback all migrations and reapply
	@echo "Resetting database..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down -all || true
	@echo "Reapplying all migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up
	@echo "✓ Database reset successfully"

migrate-force: ## Force set migration version (use: make migrate-force V=1)
	@if [ -z "$(V)" ]; then echo "Error: Please specify version with V=N"; exit 1; fi
	@echo "Forcing migration version to $(V)..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(V)
	@echo "✓ Migration version forced to $(V)"

migrate-version: ## Show current migration version
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-create: ## Create new migration files (use: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "Error: Please specify migration name with NAME=xxx"; exit 1; fi
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)
	@echo "✓ Migration files created"

db-status: ## Show database status
	@echo "=== Database Status ==="
	@echo ""
	@echo "Current migration version:"
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version || echo "No migrations applied"
	@echo ""
	@echo "Tables in database:"
	@docker exec -i car_dealership_db psql -U $(DB_USER) -d $(DB_NAME) -c "\dt"
	@echo ""
	@echo "Table row counts:"
	@docker exec -i car_dealership_db psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT 'users' as table_name, COUNT(*) as row_count FROM users UNION ALL SELECT 'cars', COUNT(*) FROM cars UNION ALL SELECT 'orders', COUNT(*) FROM orders UNION ALL SELECT 'order_items', COUNT(*) FROM order_items UNION ALL SELECT 'cart_items', COUNT(*) FROM cart_items UNION ALL SELECT 'two_factor_codes', COUNT(*) FROM two_factor_codes ORDER BY table_name;"

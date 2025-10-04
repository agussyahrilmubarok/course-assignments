##################################################
### cmd
##################################################
# Run app in development
app/dev:
	go run cmd/server/main.go

app/migrate:
	go run cmd/migrater/main.go --direction=up

##################################################
### golang-migrate
##################################################
DB_URL=postgres://demo:demopass@localhost:5432/account_db?sslmode=disable
MIGRATIONS_DIR=./migrations

# Install golang-migrate CLI
golang-migrate/install:
	@echo "📦 Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✅ golang-migrate installed successfully."

# Initialize migration folder (create directory if not exists)
golang-migrate/init:
	@echo "🗂️  Initializing migrations directory..."
	@mkdir -p $(MIGRATIONS_DIR)
	@echo "✅ Migrations folder ready at $(MIGRATIONS_DIR)"

# Create new migration files
# Usage: make migrate-new name=create_users_table
golang-migrate/new:
	@if [ -z "$(name)" ]; then echo "❌ Please provide migration name, e.g. make migrate-new name=create_users_table"; exit 1; fi
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
	@echo "✅ Migration files created for '$(name)'"

# Run all pending migrations
golang-migrate/up:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up
	@echo "✅ All migrations applied successfully."

# Rollback last migration
golang-migrate/down:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1
	@echo "↩️  Rolled back last migration."

# Reset database (all migrations down)
golang-migrate/reset:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down -all
	@echo "⚠️  Database reset (all migrations reverted)."

# Show current migration version
golang-migrate/version:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version
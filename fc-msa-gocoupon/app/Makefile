api/dev:
	@echo "running API in dev mode..."
	go run cmd/api/main.go --config config.yml

api/build:
	@echo "formatting..."
	go fmt ./...
	@echo "static analyzing..."
	go vet ./...
	@echo "building binary..."
	go build -o bin/api cmd/api/main.go
	@echo "build completed!"

#####################################################################################
### migrater
#####################################################################################
migrater/up:
	go run cmd/migrater/main.go \
	--db postgres://postgres:my-secret-pw@localhost:5432/coupon_db?sslmode=disable \
	--dir migrations \
	--action up

migrater/down:
	go run cmd/migrater/main.go \
	--db postgres://postgres:my-secret-pw@localhost:5432/coupon_db?sslmode=disable \
	--action down

# go run cmd/migrater/main.go --db <URL> --action up
# go run cmd/migrater/main.go --db <URL> --action down
# go run cmd/migrater/main.go --db <URL> --action drop
# go run cmd/migrater/main.go --db <URL> --action version
# go run cmd/migrater/main.go --db <URL> --action force --force 3
# go run cmd/migrater/main.go \
#   --db postgres://postgres:my-secret-pw@localhost:5432/coupon_db?sslmode=disable \
#   --dir migrations \
#   --action up
migrater/command:
	go run cmd/migrater/main.go \
	--db $(DB_URL) \
	--dir $(DIR) \
	--action $(C)

#####################################################################################
### swagger
#####################################################################################
swagger/install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/swaggo/swag@latest
	go mod tidy

swagger/generate:
	swag init -g internal/api/$(version)/router.go -o internal/api/$(version)/docs
backend/dev:
	go run cmd/backend/main.go

backerhub/build:
	go fmt ./...
	go vet ./...
	CGO_ENABLED=0 GOOS=linux go build -o bin/backerhub-backend cmd/backend/main.go

seeder/up:
	go run cmd/seeder/main.go

swagger/install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/swaggo/swag@latest
	go mod tidy

swagger/generate:
	swag init -g internal/rest/$(version)/router.go -o internal/rest/$(version)/docs
	
compose/postgres/up:
	docker compose -f docker-compose-postgres.yml up -d --build

compose/postgres/down:
	docker compose -f docker-compose-postgres.yml down
	docker volume prune -a
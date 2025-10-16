dev:
	go run cmd/server/main.go --config=config.yaml

build:
	go build -o build/account cmd/server/main.go

seeder:
	go run cmd/server/main.go --config=config.yaml

swagger/install:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger/generate:
	swag init -g ./cmd/server/main.go -o ./cmd/server/docs

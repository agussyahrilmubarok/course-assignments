server/dev:
	go run cmd/server/main.go --config=config.json

swagger/install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/swaggo/swag@latest
	go mod tidy

swagger/generate:
	swag init -g ./cmd/server/main.go -o ./cmd/server/docs
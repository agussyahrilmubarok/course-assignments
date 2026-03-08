server/dev:
	go run cmd/server/main.go --config=config.yaml

csv/export:
	go run cmd/csv/main.go --config=config.yaml --output=$(OUTPUT)
	
swagger/install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/swaggo/swag@latest
	go mod tidy

swagger/generate:
	swag init -g api/$(version)/handler/handler.go -o api/$(version)/docs
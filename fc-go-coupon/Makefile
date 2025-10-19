server/dev:
	go run cmd/server/main.go
	
swag/install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/swaggo/swag@latest
	go mod tidy

swag/gen:
	swag init -g api/$(version)/handler/handler.go -o api/$(version)/docs
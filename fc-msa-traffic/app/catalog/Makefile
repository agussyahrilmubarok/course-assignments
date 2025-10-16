server/dev:
	go run cmd/server/main.go --config=config.yaml

server/build:
	go build -o build/catalog cmd/server/main.go

seeder/run:
	go run cmd/server/main.go --config=config.yaml

seeder/build:
	go build -o build/seeder cmd/seeder/main.go

swagger/install:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger/generate:
	swag init -g ./cmd/server/main.go -o ./cmd/server/docs

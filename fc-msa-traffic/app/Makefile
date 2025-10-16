server/account:
	go run cmd/account/server/main.go --config=configs/account.yaml

server/catalog:
	go run cmd/account/server/main.go --config=configs/account.yaml

build/account:
	go build -o build/account ./cmd/account/server

seeder/account:
	go run cmd/account/seeder/main.go --config=configs/account.yaml

swagger/install:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger/account:
	swag init -g ./cmd/account/server/main.go -o ./cmd/account/server/docs
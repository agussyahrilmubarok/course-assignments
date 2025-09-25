app/start:
	go run cmd/main.go

compose/up:
	docker compose up -d --build

compose/down:
	docker compose down
	docker volume prune -a

swagger/install:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag
	
swagger/init:
	swag init -g cmd/main.go

swagger/url:
	http://localhost:8080/swagger/index.html
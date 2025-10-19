server/dev:
	go run cmd/server/main.go --config=backerhub-dev.json

server/stage:
	go run cmd/server/main.go --config=backerhub-stage.json

migrater/run:
	go run cmd/migrater/main.go --dsn="host=127.0.0.1 user=backeruser password=backerpass dbname=backerhub port=5432 sslmode=disable TimeZone=Asia/Jakarta"

compose/up:
	docker compose up -d --build

compose/down:
	docker compose down
	docker volume prune -a
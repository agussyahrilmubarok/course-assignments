compose/up:
	docker compose up -d --build

compose/down:
	docker compose down
	docker volume prune -a

compose/dev/up:
	docker compose -f docker-compose-dev.yml up -d --build

compose/dev/down:
	docker compose -f docker-compose-dev.yml down
	docker volume prune -a
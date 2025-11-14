compose/up:
	docker compose up -d --build

compose/down:
	docker compose down
	docker volume prune -a

compose/prod/up:
	docker compose -f docker-compose-prod.yaml up -d --build

compose/prod/down:
	docker compose -f docker-compose-prod.yaml down
	docker volume prune -a
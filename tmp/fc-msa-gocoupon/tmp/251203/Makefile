compose/up:
	docker compose up -d --build

compose/down:
	docker compose down 
	docker volume prune -a

compose/local/up:
	docker compose -f docker-compose-local.yml up -d

compose/local/down:
	docker compose -f docker-compose-local.yml down
	docker volume prune -a
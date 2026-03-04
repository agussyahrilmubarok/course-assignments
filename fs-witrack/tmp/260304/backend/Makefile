backend/dev:
	mvn clean spring-boot:run

backend/package:
	mvn clean package

backend/build:
	mvn -Pnative native:compile

compose/postgres/up:
	docker compose -f docker-compose-postgres.yml up -d --build

compose/postgres/down:
	docker compose -f docker-compose-postgres.yml down
	docker volume prune -a

compose/backend/up:
	docker compose -f docker-compose-backend.yml up -d --build

compose/backend/down:
	docker compose -f docker-compose-backend.yml down
	docker volume prune -a
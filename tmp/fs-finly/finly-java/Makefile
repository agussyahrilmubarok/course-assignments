app/start:
	mvn clean spring-boot:run

docker/compose/up:
	docker compose up -d

docker/compose/down:
	docker compose down
	docker volume prune -a

docker/image/build:
	docker build -t finly .

docker/compose/prod/up:
	docker compose -f docker-compose-prod.yml up -d

docker/compose/prod/down:
	docker compose -f docker-compose-prod.yml down
	docker volume prune -a

docker/clean:
	docker container prune && docker volume prune && docker image prune
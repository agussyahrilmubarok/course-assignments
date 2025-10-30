server/dev:
	cd $(name) && ./gradlew bootRun

server/build:
	cd $(name) && ./gradlew clean build -x test

server/run:
	java -jar $(name)/build/libs/*.jar

server/clean:
	cd $(name) && ./gradlew clean

compose/up:
	docker compose up -d --build

compose/down:
	docker compose down
	docker volume prune -a

compose/all/down:
	docker compose -f docker-compose-infra.yml down
	docker compose -f docker-compose-obs.yml down
	docker volume prune -a

firewall/allow/port:
	sudo ufw allow in on docker0
	sudo ufw allow 7071/tcp
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
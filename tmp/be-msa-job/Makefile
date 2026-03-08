loans.run:
	cd loans && mvn clean spring-boot:run -Dspring-boot.run.profiles=qa

accounts.run:
	cd accounts && mvn clean spring-boot:run -Dspring-boot.run.profiles=qa

cards.run:
	cd cards && mvn clean spring-boot:run -Dspring-boot.run.profiles=qa

configserver.run:
	cd configserver && mvn clean spring-boot:run

eurekaserver.run:
	cd eurekaserver && mvn clean spring-boot:run 

gatewayserver.run:
	cd gatewayserver && mvn clean spring-boot:run 

configserver.encrypt:
	curl -X POST http://localhost:8071/encrypt -H "Content-Type: text/plain" -d 'mysecretpassword'

docker.loans.build:
	cd loans && docker build -t examplebank/loans .

docker.accounts.build:
	cd accounts && docker build -t examplebank/accounts .

docker.cards.build:
	cd cards && docker build -t examplebank/cards .

docker.configserver.build:
	cd configserver && docker build -t examplebank/configserver .

docker.eurekaserver.build:
	cd eurekaserver && docker build -t examplebank/eurekaserver .

docker.gatewayserver.build:
	cd gatewayserver && docker build -t examplebank/gatewayserver .

docker.dev.up:
	docker compose -f docker-compose-dev.yml up -d --build

docker.dev.down:
	docker compose -f docker-compose-dev.yml down

redis.login:
	docker exec -it redis redis-cli

redis.check.db:
	docker exec redis redis-cli DBSIZE

redis.check.key:
	docker exec redis redis-cli EXISTS $(KEY)

mysql.login:
	docker exec -it mysql mysql -uroot -prootpass

mysql.check.db:
	docker exec -i mysql mysql -uroot -prootpass -e "SHOW DATABASES LIKE '$(DB)';"

mysql.check.data:
	docker exec -i mysql mysql -uroot -prootpass -e "USE $(DB); SELECT * FROM $(TABLE) LIMIT 5;"

rabbitmq.image:
	docker pull rabbitmq:management

rabbitmq.container:
	# guest
	docker run --name rabbitmq -p 5672:5672 -p 15672:15672 -d rabbitmq:management

redis.image:
	docker pull redis

redis.container:
	docker run --name redis -d redis

redis.cli:
	docker exec -it redis redis-cli
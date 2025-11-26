compose/up:
	docker compose up -d --build

compose/down:
	docker compose down
	docker volume prune -a

compose/dev/up:
	docker network create coupon_net
	docker compose -f docker-compose-infra.yml up -d
	docker compose -f docker-compose-mon.yml up -d

compose/dev/down:
	docker compose -f docker-compose-infra.yml down
	docker compose -f docker-compose-mon.yml down
	docker network rm coupon_net
	docker volume prune -a
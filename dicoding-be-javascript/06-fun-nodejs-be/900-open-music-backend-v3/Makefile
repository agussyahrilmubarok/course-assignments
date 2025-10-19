container-db:
	docker compose -f docker-compose-db.yml up

container-dev:
	docker compose -f docker-compose-dev.yml up

install:
	npm install

migrations:
	npm run migrate create $(name)

migrate:
	npm run migrate up

lint-fix:
	npx eslint ./src --fix
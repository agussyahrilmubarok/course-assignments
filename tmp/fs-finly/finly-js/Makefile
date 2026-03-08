pnpm/install:
	pnpm add bcrypt connect-flash daisyui dotenv ejs express express-session express-validator mongoose morgan pm2 uuid winston winston-daily-rotate-file
	pnpm add -D autoprefixer nodemon postcss postcss-cli prettier tailwindcss

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
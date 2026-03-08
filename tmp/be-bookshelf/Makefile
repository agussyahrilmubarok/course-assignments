################################################################################################
##### DART
################################################################################################
run:
	dart ./bin/server.dart

run-test:
	dart test

dep-installer:
	dart pub get

formatter:
	dart format --output=none --set-exit-if-changed .

prev-format:
	dart fix --dry-run

fix-format:
	dart fix --apply

analyzer:
	dart analyze

exec:
	dart compile exe ./bin/server.dart
	
################################################################################################
##### DOCKER
################################################################################################
container-db:
	docker compose -f docker-compose-db.yml up

container:
	docker compose -f docker-compose.yml up

image:
	docker build -f eg-dart-shelf .

rm-data:
	sudo rm -R postgres-data
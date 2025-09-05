server.dev:
	cd $(name) && ./gradlew bootRun

server.build:
	cd $(name) && ./gradlew clean build -x test

server.run:
	java -jar $(name)/build/libs/*.jar

server.clean:
	cd $(name) && ./gradlew clean
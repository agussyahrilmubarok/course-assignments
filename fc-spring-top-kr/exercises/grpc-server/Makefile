# ============================================
# 🧠 Spring Boot gRPC Server Makefile
# ============================================

# --- Configuration Variables ---
APP_NAME := grpc-server
JAR_FILE := target/$(APP_NAME)-0.0.1-SNAPSHOT.jar
GRPC_PORT := 9090
HTTP_PORT := 8080

# Maven Wrapper (use ./mvnw if available)
MVN := ./mvnw

# --- Commands ---

## 🧹 Clean all build artifacts
clean:
	@echo "🧹 Cleaning project..."
	$(MVN) clean

## 🧱 Compile Java & Protobuf files
compile:
	@echo "🧱 Compiling sources and generating protobuf classes..."
	$(MVN) clean compile

## 🧪 Run tests
test:
	@echo "🧪 Running tests..."
	$(MVN) test

## 🚀 Run Spring Boot gRPC server locally
run:
	@echo "🚀 Starting gRPC server on port $(GRPC_PORT) (HTTP on $(HTTP_PORT))..."
	$(MVN) spring-boot:run

## 🧩 Regenerate protobuf classes manually (if needed)
proto:
	@echo "🔧 Regenerating protobuf Java classes..."
	$(MVN) protobuf:compile
	$(MVN) protobuf:compile-custom

## 📦 Package project into JAR
package:
	@echo "📦 Packaging application..."
	$(MVN) clean package -DskipTests

## 🧰 Run packaged JAR
start:
	@echo "🚀 Running $(APP_NAME) JAR..."
	java -jar $(JAR_FILE)

## 🧱 Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(APP_NAME):latest .

## 🏃 Run Docker container
docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p $(HTTP_PORT):8080 -p $(GRPC_PORT):9090 $(APP_NAME):latest

## 🧽 Stop all Docker containers for this app
docker-stop:
	@echo "🛑 Stopping running containers..."
	docker ps -q --filter "ancestor=$(APP_NAME):latest" | xargs -r docker stop

## 🧾 Show logs of the running Spring Boot app
logs:
	@echo "📜 Showing Spring Boot logs..."
	tail -f logs/spring.log

## 💡 Help command
help:
	@echo ""
	@echo "📘 Available Makefile Commands for gRPC Server:"
	@echo "------------------------------------------------"
	@echo " make clean           - Clean project build files"
	@echo " make compile         - Compile Java and generate gRPC stubs"
	@echo " make test            - Run all tests"
	@echo " make run             - Run Spring Boot gRPC app (dev mode)"
	@echo " make proto           - Rebuild protobuf-generated classes"
	@echo " make package         - Build JAR without running tests"
	@echo " make start           - Run JAR directly"
	@echo " make docker-build    - Build Docker image"
	@echo " make docker-run      - Run app in Docker container"
	@echo " make docker-stop     - Stop Docker container"
	@echo " make logs            - Tail application logs"
	@echo " make help            - Show this help message"
	@echo ""

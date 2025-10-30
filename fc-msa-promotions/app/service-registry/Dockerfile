FROM gradle:8.10.1-jdk21 AS builder
WORKDIR /app

COPY . .
RUN gradle dependencies --no-daemon || true
RUN gradle clean bootJar -x test --no-daemon

FROM eclipse-temurin:21-jre-alpine AS runtime
WORKDIR /app

COPY --from=builder /app/build/libs/*.jar app.jar
EXPOSE 8761

ENTRYPOINT ["java", "-jar", "app.jar"]

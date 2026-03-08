FROM maven:3.9.7-amazoncorretto-21 AS build
COPY pom.xml ./
COPY .mvn .mvn
COPY src src
RUN mvn clean install -DskipTests

FROM amazoncorretto:21
WORKDIR gatewayserver
COPY --from=build target/*.jar gatewayserver.jar
EXPOSE 8072

ENTRYPOINT ["java", "-jar", "gatewayserver.jar"]
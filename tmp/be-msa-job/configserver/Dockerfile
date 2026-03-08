FROM maven:3.9.7-amazoncorretto-21 AS build
COPY pom.xml ./
COPY .mvn .mvn
COPY src src
RUN mvn clean install -DskipTests

FROM amazoncorretto:21
WORKDIR configserver
COPY --from=build target/*.jar configserver.jar
EXPOSE 8071

ENTRYPOINT ["java", "-jar", "configserver.jar"]
FROM maven:3.9.7-amazoncorretto-21 AS build
COPY pom.xml ./
COPY .mvn .mvn
COPY src src
RUN mvn clean install -DskipTests

FROM amazoncorretto:21
WORKDIR eurekaserver
COPY --from=build target/*.jar eurekaserver.jar
EXPOSE 8070

ENTRYPOINT ["java", "-jar", "eurekaserver.jar"]
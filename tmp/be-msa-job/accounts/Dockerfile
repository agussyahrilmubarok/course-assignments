FROM maven:3.9.7-amazoncorretto-21 AS build
COPY pom.xml ./
COPY .mvn .mvn
COPY src src
RUN mvn clean install -DskipTests

FROM amazoncorretto:21
WORKDIR accounts
COPY --from=build target/*.jar accounts.jar
EXPOSE 8080

ENTRYPOINT ["java", "-jar", "accounts.jar"]
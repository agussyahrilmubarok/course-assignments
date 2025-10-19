FROM maven:3.9.7-amazoncorretto-21 AS build
COPY pom.xml ./
COPY .mvn .mvn
COPY src src
RUN mvn clean install -DskipTests

FROM amazoncorretto:21
WORKDIR loans
COPY --from=build target/*.jar loans.jar
EXPOSE 8090

ENTRYPOINT ["java", "-jar", "loans.jar"]
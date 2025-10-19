FROM maven:3.9.7-amazoncorretto-21 AS build
COPY pom.xml ./
COPY .mvn .mvn
COPY src src
RUN mvn clean install -DskipTests

FROM amazoncorretto:21
WORKDIR cards
COPY --from=build target/*.jar cards.jar
EXPOSE 1133

ENTRYPOINT ["java", "-jar", "cards.jar"]
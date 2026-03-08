# demo-spring-microservices

Demo Spring Microservices Architecture for Bank API.

## How to run

## Locally

```bash
# Setup MySQL and RabbitMQ
docker compose -f docker-compose-dev.yml

# Run Config Server
make configserver.run

# Run Eureka Server
make eurekaserver.run

# Run Accounts Service
make accounts.run

# Run Cards Service
make cards.run

# Run Loans Service
make loans.run

# Open Swagger on each services and request init api
```

## Docker 

```bash
# Run Compose
docker compose up -d --build
```

## Technology

- **Spring Framework**:
  - **Spring Boot**: Rapid development of production-grade applications.
  - **Spring Cloud Config**: Centralized configuration management.
  - **Spring Cloud Gateway**: API gateway for routing requests and handling cross-cutting concerns.
  - **Spring Data JPA**: Simplifies database operations via ORM (Object Relational Mapping).
  - **Spring Cloud Netflix Eureka**: Service discovery and registration.
  - **Lombok**: Boilerplate code reduction (e.g., getters, setters, constructors).
  - **Keycloak**: Identity and access management solution for authentication and authorization (OAuth2 and OpenID Connect).
- **Resilience**:
  - **Resilience4j**: Fault tolerance with circuit breakers, rate limiters, retries, and bulkheads.
  - **Retry**: Automatic retry mechanisms for failed operations.
  - **Rate Limiting with Redis**: Redis-backed rate limiting for controlling the request flow across services.
- **Observability**:
  - **Grafana Alloy**: Metrics and log collector/agent for Prometheus, Loki, and Grafana.
  - **Grafana Loki**: Log aggregation system, optimized for storing and querying logs.
  - **Grafana**: Visualization and dashboarding tool for metrics, logs, and traces.
  - **MinIO**: S3-compatible object storage used to store log chunks and other artifacts.
  - **Micrometer**: Application metrics facade used to integrate with Prometheus and other monitoring systems.
  - **Prometheus**: Monitoring and alerting toolkit that scrapes metrics exposed by services.
- **Infrastructure**:
  - **Nginx**: Reverse proxy and load balancer for distributing incoming traffic and SSL termination.
  - **Docker & Docker Compose**: Containerization of services and easy orchestration of multi-service environments.
  - **Makefile**: Simplified and standardized command execution for build, run, and management of services.
- **Database**:
  - **MySQL**: Relational database used for persistent data storage by services.

## References

- [Spring Boot Documentation](https://docs.spring.io/spring-boot/docs/current/reference/html/)
- [Spring Cloud Netflix Eureka](https://cloud.spring.io/spring-cloud-netflix/multi/multi_spring-cloud-eureka-server.html)
- [Spring Cloud Gateway](https://docs.spring.io/spring-cloud-gateway/docs/current/reference/html/)
- [Resilience4j Documentation](https://resilience4j.readme.io/docs)
- [Grafana Loki](https://grafana.com/oss/loki/)
- [Grafana Alloy](https://grafana.com/docs/alloy/)
- [Grafana](https://grafana.com/)
- [MinIO Documentation](https://min.io/docs/minio/linux/index.html)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [MySQL Documentation](https://dev.mysql.com/doc/)
- [Keycloak Documentation](https://www.keycloak.org/documentation)

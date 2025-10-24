# Black Friday Microservices Demo

## Application Components

The application consists of several microservices, each serving a specific function in the Black Friday e-commerce ecosystem:

| Service          | Port |
| ---------------- | ---- |
| API Gateway      | 8080 |
| Service Registry | 8761 |
| Member Service   | 8081 |
| Catalog Service  | 8082 |
| Search Service   | 8083 |
| Order Service    | 8084 |
| Payment Service  | 8085 |

---

## Architecture Overview & Responsibilities

* **API Gateway**

  * Acts as a single entry point to the system.
  * Responsible for routing external requests to the appropriate internal services.

* **Service Registry**

  * Uses Eureka for service discovery.
  * Enables dynamic registration and lookup of services to ensure scalability and fault tolerance.

* **Member Service**

  * Handles user registration, authentication, and profile management.
  * Issues JWT tokens for secure communication across services.

* **Catalog Service**

  * Manages product data, including categories and tags.
  * Uses PostgreSQL for transactional data and Cassandra for scalable product catalog storage.
  * Publishes events to Kafka for downstream services (e.g., Search).

* **Search Service**

  * Consumes product tag events from Kafka.
  * Uses Redis to enable fast, in-memory keyword-based product search.

* **Order Service**

  * Handles order creation, processing, and status tracking.
  * Uses OpenFeign to communicate with Catalog (for stock checking) and Payment services.
  * Publishes order events to Kafka for audit and monitoring purposes.

* **Payment Service**

  * Integrates with external payment gateways (e.g., Midtrans) to process transactions.
  * Handles payment validation and confirmation.

---

## Technology Stack

* **Programming Language**: Java 21
* **Framework**: Spring Boot (Microservices architecture)
* **Databases**: PostgreSQL (relational), Cassandra (NoSQL), Redis (in-memory)
* **Messaging**: Apache Kafka (event-driven architecture)
* **Authentication**: JWT (JSON Web Token)
* **Service Communication**: Eureka (Service Discovery), OpenFeign (HTTP client)
* **External Integration**: Midtrans (Payment Gateway)

---

## Conclusion

This microservices-based system demonstrates a scalable and modular architecture tailored for high-traffic events like Black Friday. By leveraging modern technologies such as Kafka for event-driven communication, Redis for fast data access, and service discovery via Eureka, the system ensures high performance, scalability, and resilience. Each microservice is designed to handle specific business functions, allowing for independent development, deployment, and scaling.

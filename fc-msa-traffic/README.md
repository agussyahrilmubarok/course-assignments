# Demo: Traffic Control in Microservices with Golang

## Application Components

This demo consists of the following microservices and supporting infrastructure:

* **API Gateway** – Port `8080`
* **Account Service** – Port `8081`
* **Catalog Service** – Port `8082`
* **Order Service** – Port `8083`
* **Pricing Service** – Port `8084`
* **PostgreSQL Database** – Port `5432`
* **Redis** – Port `6379`
* **Consul Service Discovery** – Port `8500`

## Insights

* Account Service

  * Manages **authentication and authorization** using **JWT (JSON Web Tokens)**.
  * Implements a **Rate Limiter** using the **Echo Framework** to control traffic and prevent abuse.

* Catalog Service

  * Handles **product management** for the e-commerce platform.
  * Applies a **seeder mechanism** to initialize product data during startup.
  * Provides **Inter-process Communication (IPC)** to coordinate and control product stock levels.
  * Implements **caching** to accelerate access to frequently requested stock data.

* Order Service

  * Manages **product order creation** workflows.
  * Utilizes **Goroutines** for concurrent operations, including:

    * **Stock validation** via communication with the Catalog Service.
    * **Dynamic pricing retrieval** from the Pricing Service.
  * Implements Database Sharding to distribute data across multiple databases for scalability and performance optimization, using examples such as PostgreSQL and MySQL.

* Pricing Service

  * Manages **dynamic product pricing**, particularly for **time-based promotions** such as **flash sales**.
  * Provides pricing logic and adjustment mechanisms based on promotional rules.

## Conclusion

This demo illustrates how traffic control mechanisms such as rate limiting, caching, and asynchronous communication can be effectively implemented in a Golang-based microservice architecture. The system demonstrates scalable service orchestration using tools like Consul, Redis, PostgreSQL and MySQL, and leverages concurrency features of Golang to improve performance and responsiveness. This setup serves as a solid foundation for building robust, distributed e-commerce platforms.

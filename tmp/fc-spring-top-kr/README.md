# Spring Example

This project demonstrates multiple communication technologies within a Spring-based application, including REST, gRPC, GraphQL, and WebSocket. It is intended as a practical reference and learning resource for modern API development using the Spring ecosystem.

## Application Features

* **REST API** – Traditional HTTP-based request/response interactions
* **gRPC (server & client)** – High-performance binary communication using Protocol Buffers
* **GraphQL** – Flexible query-driven data fetching
* **WebSocket** – Real-time bidirectional communication

## Learning Objectives / Insights

This project helps developers:

* Understand REST API principles and implementation in Spring
* Learn how to build and consume gRPC services
* Explore data querying and schema design with GraphQL
* Implement event-driven messaging using WebSocket

## Technology Stack

* **Java** (JDK 21 or later recommended)
* **Spring Framework / Spring Boot**
* **gRPC & Protobuf**
* **GraphQL Java / Spring GraphQL**
* **WebSocket**

## API Endpoints & Usage

Each communication method may provide different endpoints or flows:

* REST endpoints available under:
  `http://localhost:<port>/api/...`

* gRPC server runs on configured port in application settings

* GraphQL playground or endpoint typically at:
  `http://localhost:<port>/graphql`

* WebSocket connection through:
  `ws://localhost:<port>/ws/...`

(Detailed endpoint documentation may be added as the project evolves.)

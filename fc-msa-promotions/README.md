# Demo Microservice Architecture in promotion feature written by Java with Spring Framework

# Applications

* service-registry,8761
* user-service,8081,7071
* coupon-service,8082,7072
* postgres,5432
* redis,6379

# Insight

* service-registry
    * The purpose is to avoid hardcoding service addresses. With Eureka Server, each service can register itself and discover other services dynamically.
    * Eureka works with client-side load balancing to distribute requests across multiple service instances.
    * If one instance goes down, Eureka automatically removes it from the registry so requests won’t be routed to a dead service.
* user-service
    * Using Spring Security + JWT (JSON Web Token), users can log in, and all subsequent requests to other services use the token.
    * In a microservice architecture, only the user service should be the authoritative source for user identity.
    * The user service is accessed via the API Gateway for login/registration, and the gateway forwards the token to other services.
* coupon-service
    * Manage coupon creation, distribution, and usage based on Coupon Policy and coupon status.
    * With @Transactional, operations such as issue coupon, use coupon, or cancel coupon can be executed consistently. If an error occurs midway, the database transaction will roll back to avoid data corruption.
    * Handle validation rules like expired coupons, already used coupons, or coupons outside the valid time frame.
    * The coupon service focuses only on coupon management and does not store user data. User identity comes from a JWT token or through a UserIdInterceptor.
* api-gateway
    * The purpose is to provide a single entry point to all microservices. Clients (mobile/web) don’t need to know the addresses of internal services.
    * For example, requests to /users/** go to user-service, and /coupons/** go to coupon-service. Routing is flexible and can be configured dynamically.
    * The gateway can have pre-filters to validate JWT tokens, check user roles, or log requests. This way, internal services don’t have to repeat authentication logic.
    * The API Gateway can also handle rate limiting, circuit breaking, monitoring, and logging.
* 
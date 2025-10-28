package com.example.gateway.filter;

import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.springframework.cloud.client.loadbalancer.reactive.ReactorLoadBalancerExchangeFilterFunction;
import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.cloud.gateway.filter.factory.AbstractGatewayFilterFactory;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

import java.util.Map;

@Component
@Slf4j
public class JwtAuthenticationFilter extends AbstractGatewayFilterFactory<JwtAuthenticationFilter.Config> {

    private final WebClient webClient;

    public JwtAuthenticationFilter(ReactorLoadBalancerExchangeFilterFunction lbFunction) {
        super(Config.class);
        this.webClient = WebClient.builder()
                .filter(lbFunction)
                .baseUrl("http://USER-SERVICE")
                .build();
    }

    @Override
    public GatewayFilter apply(Config config) {
        return (exchange, chain) -> {
            String authHeader = exchange.getRequest().getHeaders().getFirst("Authorization");
            if (authHeader != null && authHeader.startsWith("Bearer ")) {
                String token = authHeader.substring(7);
                return validateToken(token)
                        .flatMap(userId -> proceedWithUserId(userId, exchange, chain))
                        .switchIfEmpty(chain.filter(exchange)) // If token is invalid, continue without setting userId
                        .onErrorResume(e -> handleAuthenticationError(exchange, e)); // Handle errors
            }

            return chain.filter(exchange);
        };
    }

    private Mono<Void> handleAuthenticationError(ServerWebExchange exchange, Throwable e) {
        exchange.getResponse().setStatusCode(HttpStatus.UNAUTHORIZED);
        return exchange.getResponse().setComplete();
    }

    private Mono<String> validateToken(String token) {
        log.info("Validating token: {}", token);
        return webClient.post()
                .uri("/api/v1/auth/validate-token")
                .bodyValue("{\"token\":\"" + token + "\"}")
                .header("Content-Type", "application/json")
                .retrieve()
                .bodyToMono(Map.class)
                .flatMap(response -> {
                    log.info("Received response from token validation: {}", response);
                    Object userObj = response.get("user");
                    if (userObj instanceof Map<?, ?> userMap) {
                        Object idObj = userMap.get("id");
                        if (idObj != null) {
                            log.info("Extracted userId from response: {}", idObj);
                            return Mono.just(idObj.toString());
                        }
                    }
                    return Mono.empty();
                });
    }

    private Mono<Void> proceedWithUserId(String userId, ServerWebExchange exchange, GatewayFilterChain chain) {
        return chain.filter(
                exchange.mutate()
                        .request(exchange.getRequest().mutate()
                                .header("X-USER-ID", userId)
                                .build()
                        )
                        .build()
        );
    }

    public static class Config {
        // Configuration class for filter setup
    }

    @Data
    public static class TokenValidationResponse {
        private Long id;
        private String email;
        private Boolean valid;
    }
}

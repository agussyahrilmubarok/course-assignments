package com.example.gateway.config;

import org.springframework.cloud.gateway.filter.ratelimit.KeyResolver;
import org.springframework.cloud.gateway.filter.ratelimit.RedisRateLimiter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import reactor.core.publisher.Mono;

@Configuration
public class RateLimitConfig {

    @Bean
    public RedisRateLimiter redisRateLimiter() {
        // replenishRate: Number of requests allowed per second
        // burstCapacity: Maximum accumulated requests allowed
        return new RedisRateLimiter(10, 20);
    }

    @Bean
    public KeyResolver userKeyResolver() {
        return exchange -> Mono.just(
                exchange.getRequest().getHeaders().getFirst("X-USER-ID") != null ?
                        exchange.getRequest().getHeaders().getFirst("X-USER-ID") :
                        exchange.getRequest().getRemoteAddress().getAddress().getHostAddress()
        );
    }
}

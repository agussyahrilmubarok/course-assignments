package com.example.gateway.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Mono;

import java.util.Map;

@RestController
@RequestMapping("/fallback")
public class FallbackController {

    @GetMapping("/accounts")
    public Mono<Map<String, Object>> accountFallback() {
        return Mono.just(Map.of("status", "down"));
    }

    @GetMapping("/catalogs")
    public Mono<Map<String, Object>> catalogFallback() {
        return Mono.just(Map.of("status", "down"));
    }

    @GetMapping("/orders")
    public Mono<Map<String, Object>> orderFallback() {
        return Mono.just(Map.of("status", "down"));
    }

    @GetMapping("/pricings")
    public Mono<Map<String, Object>> pricingFallback() {
        return Mono.just(Map.of("status", "down"));
    }
}

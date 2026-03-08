package com.example.gateway.config;

import io.micrometer.observation.ObservationRegistry;
import org.springframework.cloud.client.loadbalancer.LoadBalanced;
import org.springframework.cloud.client.loadbalancer.reactive.ReactorLoadBalancerExchangeFilterFunction;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.reactive.function.client.WebClient;

@Configuration
public class WebClientConfig {

    @Bean
    @LoadBalanced
    public WebClient webClient(
            ReactorLoadBalancerExchangeFilterFunction lbFunction,
            ObservationRegistry observationRegistry) {
        return WebClient.builder()
                .filter(lbFunction)
                .observationRegistry(observationRegistry)
                .build();
    }
}

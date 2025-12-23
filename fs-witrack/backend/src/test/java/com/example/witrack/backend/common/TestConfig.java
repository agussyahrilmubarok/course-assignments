package com.example.witrack.backend.common;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;

@TestConfiguration
class TestConfig {

    @Bean
    ObjectMapper objectMapper() {
        return new ObjectMapper();
    }
}

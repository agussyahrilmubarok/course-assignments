package com.example.order.config;

import jakarta.servlet.SessionTrackingMode;
import java.util.Set;
import org.springframework.boot.web.servlet.ServletContextInitializer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;


@Configuration
public class ServletConfig {

    @Bean
    public ServletContextInitializer servletContextInitializer() {
        // don't append the session id to resources
        return servletContext -> servletContext.setSessionTrackingModes(Set.of(SessionTrackingMode.COOKIE));
    }

}

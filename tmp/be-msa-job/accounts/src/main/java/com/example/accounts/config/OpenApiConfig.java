package com.example.accounts.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenApiConfig {

    public OpenAPI customOpenAPI() {
        return new OpenAPI()
                .info(new Info()
                        .title("Accounts Microservices API")
                        .version("1.0.0")
                        .description("This is the OpenAPI documentation for the E-Banking system, accounts microservices."));
    }

}

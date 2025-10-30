package com.example.coupon.config;

import io.swagger.v3.oas.annotations.enums.ParameterIn;
import io.swagger.v3.oas.models.Components;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.media.*;
import io.swagger.v3.oas.models.parameters.Parameter;
import org.springdoc.core.models.GroupedOpenApi;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class SwaggerConfig {

    @Bean
    public OpenAPI openApiSpec() {
        return new OpenAPI()
                .info(new Info()
                        .title("Coupon Service API in Promotions")
                        .version("1.0.0"))
                .components(new Components()
                        .addSchemas("ApiErrorResponse", new ObjectSchema()
                                .addProperty("status", new IntegerSchema())
                                .addProperty("code", new StringSchema())
                                .addProperty("message", new StringSchema())
                                .addProperty("fieldErrors", new ArraySchema().items(
                                        new Schema<ArraySchema>().$ref("ApiFieldError"))))
                        .addSchemas("ApiFieldError", new ObjectSchema()
                                .addProperty("code", new StringSchema())
                                .addProperty("message", new StringSchema())
                                .addProperty("property", new StringSchema())
                                .addProperty("rejectedValue", new ObjectSchema())
                                .addProperty("path", new StringSchema()))
                        .addParameters("X-USER-ID", new Parameter()
                                .in(ParameterIn.HEADER.toString())
                                .name("X-USER-ID")
                                .required(true)
                                .description("User ID header")
                                .schema(new StringSchema())));
    }

    @Bean
    public GroupedOpenApi openApiGroupSpec() {
        return GroupedOpenApi.builder()
                .group("basic")
                .packagesToScan("com.example.coupon.rest.v0")
                .build();
    }

    @Bean
    public GroupedOpenApi openApiGroupV1Spec() {
        return GroupedOpenApi.builder()
                .group("v1")
                .packagesToScan("com.example.coupon.rest.v1")
                .build();
    }

    @Bean
    public GroupedOpenApi openApiGroupV2Spec() {
        return GroupedOpenApi.builder()
                .group("v2")
                .packagesToScan("com.example.coupon.rest.v2")
                .build();
    }

    @Bean
    public GroupedOpenApi openApiGroupV3Spec() {
        return GroupedOpenApi.builder()
                .group("v3")
                .packagesToScan("com.example.coupon.rest.v3")
                .build();
    }
}

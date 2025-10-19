package com.example.catalog.model;

import com.example.catalog.cassandra.domain.Product;
import com.fasterxml.jackson.annotation.JsonInclude;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.*;

import java.util.List;

@Getter
@Setter
public class ProductDTO {

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class RegisterRequest {

        @NotBlank(message = "Seller ID must not be blank")
        private String sellerId;

        @NotBlank(message = "Product name must not be blank")
        private String name;

        private String description;

        @NotNull(message = "Price is required")
        @Min(value = 0, message = "Price must be >= 0")
        private Long price;

        @NotNull(message = "Stock count is required")
        @Min(value = 0, message = "Stock count must be >= 0")
        private Long stockCount;

        private List<String> tags;
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class DecreaseStockRequest {

        @NotNull(message = "Decrease count is required")
        @Min(value = 1, message = "Decrease count must be >= 1")
        private Long decreaseCount;
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class ProductTagsMessage {

        private String productId;

        private List<String> tags;
    }

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class Response {

        private String id;
        private String sellerId;
        private String name;
        private String description;
        private Long price;
        private Long stockCount;
        private List<String> tags;

        public static Response from(Product product) {
            return Response.builder()
                    .id(product.getId())
                    .sellerId(product.getSellerId())
                    .name(product.getName())
                    .description(product.getDescription())
                    .price(product.getPrice())
                    .stockCount(product.getStockCount())
                    .tags(product.getTags())
                    .build();
        }
    }
}


package com.example.order.model;

import com.example.order.domain.ProductOrder;
import com.fasterxml.jackson.annotation.JsonInclude;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.*;

import java.time.LocalDateTime;

@Getter
@Setter
public class OrderDTO {

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class StartOrderRequest {

        @NotBlank(message = "User id cannot be blank")
        private String userId;

        @NotBlank(message = "Product id cannot be blank")
        private String productId;

        @NotNull(message = "Product count cannot be null")
        @Positive(message = "Product count must be greater than 0")
        private Long count;
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class FinishOrderRequest {

        @NotBlank(message = "Order id cannot be blank")
        private String orderId;

        @NotBlank(message = "Payment id cannot be blank")
        private String paymentId;
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class DecreaseStockCountRequest {

        @NotNull(message = "Decrease count cannot be null")
        @Positive(message = "Decrease count must be greater than 0")
        private Long decreaseCount;
    }

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class StartOrderResponse {
        private String orderId;
        private String paymentUrl;
    }

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class Response {

        private String id;
        private String userId;
        private String productId;
        private String paymentId;
        private String orderStatus;
        private LocalDateTime createdAt;
        private LocalDateTime updatedAt;

        public static Response from(ProductOrder productOrder) {
            return Response.builder()
                    .id(productOrder.getProductId())
                    .userId(productOrder.getUserId())
                    .productId(productOrder.getProductId())
                    .paymentId(productOrder.getPaymentId())
                    .orderStatus(productOrder.getOrderStatus().toString())
                    .createdAt(productOrder.getCreatedAt())
                    .updatedAt(productOrder.getUpdatedAt())
                    .build();
        }
    }
}

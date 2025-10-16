package com.example.order.model;

import com.fasterxml.jackson.annotation.JsonInclude;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Positive;
import lombok.*;

import java.time.LocalDateTime;
import java.util.Map;

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

        @NotBlank(message = "Delivery id cannot be blank")
        private String deliveryId;
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
        private Map<String, Object> paymentMethod;
        private Map<String, Object> address;
    }

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class ProductOrder {

        private String id;
        private String userId;
        private String productId;
        private String paymentId;
        private String deliveryId;
        private String orderStatus;
        private LocalDateTime createdAt;
        private LocalDateTime updatedAt;

        public static ProductOrder from(ProductOrder productOrder) {
            return ProductOrder.builder()
                    .id(productOrder.getProductId())
                    .userId(productOrder.getUserId())
                    .productId(productOrder.getProductId())
                    .paymentId(productOrder.getPaymentId())
                    .deliveryId(productOrder.getDeliveryId())
                    .orderStatus(productOrder.getOrderStatus())
                    .createdAt(productOrder.getCreatedAt())
                    .updatedAt(productOrder.getUpdatedAt())
                    .build();
        }
    }

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class ProductDetail {

        private String orderId;
        private String userId;
        private String productId;
        private String paymentId;
        private String deliveryId;
        private String orderStatus;

        @Pattern(regexp = "PENDING|PAID|FAILED",
                message = "Invalid payment status")
        private String paymentStatus;

        @Pattern(regexp = "PENDING|SHIPPED|DELIVERED|RETURNED",
                message = "Invalid delivery status")
        private String deliveryStatus;

        public static ProductDetail from(ProductOrder productOrder, String paymentStatus, String deliveryStatus) {
            return ProductDetail.builder()
                    .orderId(productOrder.getProductId())
                    .userId(productOrder.getUserId())
                    .productId(productOrder.getProductId())
                    .paymentId(productOrder.getPaymentId())
                    .deliveryId(productOrder.getDeliveryId())
                    .orderStatus(productOrder.getOrderStatus())
                    .paymentStatus(paymentStatus)
                    .deliveryStatus(deliveryStatus)
                    .build();
        }
    }
}

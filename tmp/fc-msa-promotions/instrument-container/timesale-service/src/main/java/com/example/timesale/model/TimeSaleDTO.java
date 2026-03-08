package com.example.timesale.model;

import com.example.timesale.domain.TimeSale;
import jakarta.validation.Valid;
import jakarta.validation.constraints.*;
import lombok.*;

import java.time.LocalDateTime;

@Getter
@Setter
public class TimeSaleDTO {

    @Getter
    @Setter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ProductRequest {
        @NotBlank(message = "Product name is required.")
        private String name;

        @NotNull(message = "Price is required.")
        @Positive(message = "Price must be positive.")
        private Long price;
    }

    @Getter
    @Setter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class CreateRequest {

        @Valid
        @NotNull(message = "Product information is required.")
        private ProductRequest product;

        @NotNull(message = "Quantity is required.")
        @Positive(message = "Quantity must be positive.")
        private Long quantity;

        @NotNull(message = "Discount price is required.")
        @Positive(message = "Discount price must be positive.")
        private Long discountPrice;

        @NotNull(message = "Start time is required.")
        // @FutureOrPresent(message = "Start time must be current time or in the future.") // Uncomment for production
        private LocalDateTime startAt;

        @NotNull(message = "End time is required.")
        @Future(message = "End time must be in the future.")
        private LocalDateTime endAt;
    }

    @Getter
    @Setter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class PurchaseRequest {
        @NotBlank(message = "Time sale id must not be blank.")
        private String timeSaleId;

        @NotNull(message = "Quantity must not be null.")
        @Min(value = 1, message = "Quantity must be greater than 0.")
        private Long quantity;
    }

    @Getter
    @Setter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class PurchaseRequestMessage {
        private String timeSaleId;
        private String userId;
        private Long quantity;
        private String requestId;
    }

    @Getter
    @Builder
    public static class Response {
        private String id;
        private String productId;
        private Long quantity;
        private Long remainingQuantity;
        private Long discountPrice;
        private LocalDateTime startAt;
        private LocalDateTime endAt;
        private LocalDateTime createdAt;
        private String status;

        public static Response from(TimeSale timeSale) {
            return Response.builder()
                    .id(timeSale.getId())
                    .productId(timeSale.getProduct().getId())
                    .quantity(timeSale.getQuantity())
                    .remainingQuantity(timeSale.getRemainingQuantity())
                    .discountPrice(timeSale.getDiscountPrice())
                    .startAt(timeSale.getStartAt())
                    .endAt(timeSale.getEndAt())
                    .createdAt(timeSale.getCreatedAt())
                    .status(timeSale.getStatus().name())
                    .build();
        }
    }

    @Getter
    @Builder
    public static class PurchaseResponse {
        private String timeSaleId;
        private String userId;
        private String productId;
        private Long quantity;
        private Long discountPrice;
        private LocalDateTime purchasedAt;
        private Long totalWaiting;

        public static PurchaseResponse from(TimeSale timeSale, String userId, Long quantity) {
            return PurchaseResponse.builder()
                    .timeSaleId(timeSale.getId())
                    .userId(userId)
                    .productId(timeSale.getProduct().getId())
                    .quantity(quantity)
                    .discountPrice(timeSale.getDiscountPrice())
                    .purchasedAt(LocalDateTime.now())
                    .build();
        }
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class AsyncPurchaseResponse {
        private String requestId;
        private String status;
        private Integer queuePosition;
        private Long totalWaiting;
    }
}

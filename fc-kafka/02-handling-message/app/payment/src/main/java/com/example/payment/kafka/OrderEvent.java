package com.example.payment.kafka;

import lombok.*;

import java.math.BigDecimal;
import java.time.OffsetDateTime;
import java.util.UUID;

@Getter
@Setter
@Builder
public class OrderEvent {
    private UUID eventId;              // unique per event (important for idempotency)
    private OrderEventType eventType;  // CREATED, CANCELLED
    private OffsetDateTime eventAt;    // when the event was published
    private OrderPayload order;        // order details related to this event

    public enum OrderEventType {
        CREATED,
        PAID,
        PROCESSING,
        SHIPPED,
        DELIVERED,
        CANCELLED
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class OrderPayload {
        private String orderId;
        private String customerId;
        private BigDecimal totalAmount;
        private OffsetDateTime orderAt;
    }
}

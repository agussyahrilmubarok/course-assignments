package com.example.order.domain;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Entity
@Table(name = "ProductOrders")
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
public class ProductOrder {

    @Id
    @Column(nullable = false, updatable = false)
    private String id;

    @Column(nullable = false)
    private String userId;

    @Column(nullable = false)
    private String productId;

    @Column(nullable = false)
    private Long count;

    @Column(nullable = false)
    private BigDecimal amount;

    @Column
    private String paymentId;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 32)
    private Status orderStatus = Status.CREATED;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private LocalDateTime updatedAt;

    public enum Status {
        CREATED,
        PENDING_PAYMENT,
        PAYMENT_CONFIRMED,
        PROCESSING,
        SHIPPED,
        DELIVERED,
        CANCELLED,
        REJECTED,
        RETURN_REQUESTED,
        RETURNED,
        REFUNDED
    }
}

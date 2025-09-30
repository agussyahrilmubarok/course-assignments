package com.example.payment.domain;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.hibernate.annotations.UuidGenerator;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.math.BigDecimal;
import java.time.OffsetDateTime;
import java.util.UUID;


@Entity
@Table(name = "Payments")
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
public class Payment {

    @Id
    @Column(nullable = false, updatable = false, columnDefinition = "UUID")
    @GeneratedValue
    @UuidGenerator
    private UUID id;

    @Column(nullable = false)
    private String orderId;

    @Column(nullable = false)
    private String customerId;

    @Column(nullable = false, precision = 10, scale = 2)
    private BigDecimal amount;

    @Column
    private String method;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 50)
    private Status status;

    @Column
    private OffsetDateTime paidAt;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private OffsetDateTime dateCreated;

    @LastModifiedDate
    @Column(nullable = false)
    private OffsetDateTime lastUpdated;

    public enum Status {
        PENDING,     // created but not yet processed
        PROCESSING,  // payment is being processed
        PAID,        // payment completed successfully
        FAILED,      // payment attempt failed
        CANCELLED    // payment was cancelled
    }
}

package com.example.timesale.domain;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "TimeSaleOrders")
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
public class TimeSaleOrder {

    @Id
    @Column(nullable = false, updatable = false)
    private String id;

    @Column(nullable = false)
    private String userId;

    @Column(nullable = false)
    private Long quantity;

    @Column(nullable = false)
    private Long discountPrice;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private TimeSaleOrder.Status status;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "time_sale_id", nullable = false)
    private TimeSale timeSale;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private LocalDateTime updatedAt;

    public void complete() {
        this.status = TimeSaleOrder.Status.COMPLETED;
    }

    public void fail() {
        this.status = TimeSaleOrder.Status.FAILED;
    }

    public enum Status {
        PENDING,
        COMPLETED,
        FAILED
    }
}

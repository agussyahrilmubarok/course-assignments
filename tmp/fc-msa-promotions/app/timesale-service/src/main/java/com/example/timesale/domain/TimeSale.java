package com.example.timesale.domain;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.hibernate.proxy.HibernateProxy;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "TimeSales")
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
public class TimeSale {

    @Id
    @Column(nullable = false, updatable = false)
    private String id;

    @Column(nullable = false)
    private Long quantity;

    @Column(nullable = false)
    private Long remainingQuantity;

    @Column(nullable = false)
    private Long discountPrice;

    @Column(nullable = false)
    private LocalDateTime startAt;

    @Column(nullable = false)
    private LocalDateTime endAt;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private TimeSale.Status status = TimeSale.Status.ACTIVE;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "product_id", nullable = false)
    private Product product;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private LocalDateTime updatedAt;

    public boolean isActive() {
        return status == TimeSale.Status.ACTIVE;
    }

    public void purchase(Long quantity) {
        validatePurchase(quantity);
        this.remainingQuantity -= quantity;
    }

    private void validatePurchase(Long quantity) {
        validateStatus();
        validateQuantity(quantity);
        validatePeriod();
    }

    private void validateStatus() {
        if (status != TimeSale.Status.ACTIVE) {
            throw new IllegalStateException("Time sale is not active");
        }
    }

    private void validateQuantity(Long quantity) {
        if (remainingQuantity < quantity) {
            throw new IllegalStateException("Not enough quantity available");
        }
    }

    private void validatePeriod() {
        LocalDateTime now = LocalDateTime.now();
        if (now.isBefore(startAt) || now.isAfter(endAt)) {
            throw new IllegalStateException("Time sale is not in valid period");
        }
    }

    public Product getProduct() {
        if (this.product instanceof HibernateProxy) {
            return (Product) ((HibernateProxy) this.product).getHibernateLazyInitializer().getImplementation();
        }
        return this.product;
    }

    public enum Status {
        ACTIVE,
        SOLD_OUT,
        ENDED
    }
}

package com.example.coupon.domain;

import com.example.coupon.exception.CouponAlreadyUsedException;
import com.example.coupon.exception.CouponExpiredException;
import com.example.coupon.exception.CouponNotUsedException;
import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "Coupons")
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
public class Coupon {

    @Id
    @Column(nullable = false, updatable = false)
    private String id;

    @Column(nullable = false, unique = true)
    private String code;

    @Enumerated(EnumType.STRING)
    private Status status;

    @Column
    private LocalDateTime usedAt;

    @Column(nullable = false)
    private String userId;

    @Column
    private String orderId;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "coupon_policy_id", nullable = false)
    private CouponPolicy couponPolicy;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private LocalDateTime updatedAt;

    public boolean isExpired() {
        LocalDateTime now = LocalDateTime.now();
        return now.isBefore(couponPolicy.getStartTime()) || now.isAfter(couponPolicy.getEndTime());
    }

    public boolean isUsed() {
        return status == Status.USED;
    }

    public void use(String orderId) {
        if (status == Status.USED) {
            throw new CouponAlreadyUsedException("This coupon has already been used.");
        }
        if (isExpired()) {
            throw new CouponExpiredException("This coupon has expired.");
        }
        this.status = Status.USED;
        this.orderId = orderId;
        this.usedAt = LocalDateTime.now();
    }

    public void cancel() {
        if (status != Status.USED) {
            throw new CouponNotUsedException("This coupon has not been used.");
        }
        this.status = Status.CANCELED;
        this.orderId = null;
        this.usedAt = null;
    }

    public enum Status {
        AVAILABLE,
        USED,
        EXPIRED,
        CANCELED
    }
}

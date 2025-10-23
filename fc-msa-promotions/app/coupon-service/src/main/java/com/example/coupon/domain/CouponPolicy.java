package com.example.coupon.domain;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;
import java.util.Set;

@Entity
@Table(name = "CouponPolicies")
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
public class CouponPolicy {

    @Id
    @Column(nullable = false, updatable = false)
    private String id;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false, columnDefinition = "text")
    private String description;

    @Enumerated(EnumType.STRING)
    private DiscountType discountType;

    @Column(nullable = false)
    private Integer discountValue;

    @Column(nullable = false)
    private Integer minimumOrderAmount;

    @Column(nullable = false)
    private Integer maximumDiscountAmount;

    @Column(nullable = false)
    private Integer totalQuantity;

    @Column(nullable = false)
    private LocalDateTime startTime;

    @Column(nullable = false)
    private LocalDateTime endTime;

    @OneToMany(mappedBy = "couponPolicy", fetch = FetchType.EAGER)
    private Set<Coupon> coupons;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private LocalDateTime updatedAt;

    public boolean isValidPeriod() {
        LocalDateTime now = LocalDateTime.now();
        return !now.isBefore(startTime) && !now.isAfter(endTime);
    }

    public Integer getIssuedQuantity() {
        return coupons != null ? coupons.size() : 0;
    }

    public enum DiscountType {
        FIXED_AMOUNT,
        PERCENTAGE
    }
}

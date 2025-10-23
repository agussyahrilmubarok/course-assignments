package com.example.coupon.model;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.fasterxml.jackson.annotation.JsonInclude;
import jakarta.validation.constraints.NotBlank;
import lombok.*;

import java.time.LocalDateTime;

@Getter
@Setter
public class CouponDTO {

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class IssueRequest {
        @NotBlank(message = "Coupon policy id is required.")
        private String couponPolicyId;
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class IssueMessage {

        @NotBlank(message = "Coupon policy id is required.")
        private String couponPolicyId;

        private String userId;
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class UseRequest {
        @NotBlank(message = "Order id is required.")
        private String orderId;
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class ListRequest {
        private Coupon.Status status;
        private Integer page;
        private Integer size;
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class Response {
        private String id;
        private String userId;
        private String orderId;
        private String couponCode;
        private CouponPolicy.DiscountType discountType;
        private int discountValue;
        private int minimumOrderAmount;
        private int maximumDiscountAmount;
        private LocalDateTime validFrom;
        private LocalDateTime validUntil;
        private Coupon.Status status;
        private LocalDateTime usedAt;

        public static CouponDTO.Response from(Coupon coupon) {
            CouponPolicy policy = coupon.getCouponPolicy();
            return CouponDTO.Response.builder()
                    .id(coupon.getId())
                    .userId(coupon.getUserId())
                    .couponCode(coupon.getCode())
                    .discountType(policy.getDiscountType())
                    .discountValue(policy.getDiscountValue())
                    .minimumOrderAmount(policy.getMinimumOrderAmount())
                    .maximumDiscountAmount(policy.getMaximumDiscountAmount())
                    .validFrom(policy.getStartTime())
                    .validUntil(policy.getEndTime())
                    .status(coupon.getStatus())
                    .orderId(coupon.getOrderId())
                    .usedAt(coupon.getUsedAt())
                    .build();
        }
    }
}

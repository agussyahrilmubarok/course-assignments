package com.example.coupon.model;

import com.example.coupon.domain.CouponPolicy;
import com.fasterxml.jackson.annotation.JsonInclude;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.*;

import java.time.LocalDateTime;
import java.util.UUID;

@Getter
@Setter
public class CouponPolicyDTO {

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class CreateRequest {

        @NotBlank(message = "Coupon policy name is required.")
        private String name;

        @NotBlank(message = "Coupon policy description is required.")
        private String description;

        @NotNull(message = "Discount type is required.")
        private CouponPolicy.DiscountType discountType;

        @NotNull(message = "Discount value is required.")
        @Min(value = 1, message = "Discount value must be at least 1.")
        private Integer discountValue;

        @NotNull(message = "Minimum order amount is required.")
        @Min(value = 0, message = "Minimum order amount must be at least 0.")
        private Integer minimumOrderAmount;

        @NotNull(message = "Maximum discount amount is required.")
        @Min(value = 1, message = "Maximum discount amount must be at least 1.")
        private Integer maximumDiscountAmount;

        @NotNull(message = "Total quantity is required.")
        @Min(value = 1, message = "Total quantity must be at least 1.")
        private Integer totalQuantity;

        @NotNull(message = "Start time is required.")
        private LocalDateTime startTime;

        @NotNull(message = "End time is required.")
        private LocalDateTime endTime;

        public CouponPolicy toEntity() {
            CouponPolicy couponPolicy = new CouponPolicy();
            couponPolicy.setId(UUID.randomUUID().toString());
            couponPolicy.setName(name);
            couponPolicy.setDescription(description);
            couponPolicy.setDiscountType(discountType);
            couponPolicy.setDiscountValue(discountValue);
            couponPolicy.setMinimumOrderAmount(minimumOrderAmount);
            couponPolicy.setMaximumDiscountAmount(maximumDiscountAmount);
            couponPolicy.setTotalQuantity(totalQuantity);
            couponPolicy.setStartTime(startTime);
            couponPolicy.setEndTime(endTime);

            return couponPolicy;
        }
    }

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class Response {

        private String id;
        private String name;
        private String description;
        private CouponPolicy.DiscountType discountType;
        private Integer discountValue;
        private Integer minimumOrderAmount;
        private Integer maximumDiscountAmount;
        private Integer totalQuantity;
        private Integer issuedQuantity;
        private LocalDateTime startTime;
        private LocalDateTime endTime;
        private LocalDateTime createdAt;
        private LocalDateTime updatedAt;

        public static CouponPolicyDTO.Response from(CouponPolicy couponPolicy) {
            return CouponPolicyDTO.Response.builder()
                    .id(couponPolicy.getId())
                    .name(couponPolicy.getName())
                    .description(couponPolicy.getDescription())
                    .discountType(couponPolicy.getDiscountType())
                    .discountValue(couponPolicy.getDiscountValue())
                    .minimumOrderAmount(couponPolicy.getMinimumOrderAmount())
                    .maximumDiscountAmount(couponPolicy.getMaximumDiscountAmount())
                    .totalQuantity(couponPolicy.getTotalQuantity())
                    .issuedQuantity(couponPolicy.getIssuedQuantity())
                    .startTime(couponPolicy.getStartTime())
                    .endTime(couponPolicy.getEndTime())
                    .createdAt(couponPolicy.getCreatedAt())
                    .updatedAt(couponPolicy.getUpdatedAt())
                    .build();
        }
    }
}

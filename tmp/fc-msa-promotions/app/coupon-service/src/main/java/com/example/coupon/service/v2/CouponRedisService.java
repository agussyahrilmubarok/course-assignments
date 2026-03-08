package com.example.coupon.service.v2;

import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.model.CouponDTO;

import java.util.Map;

public interface CouponRedisService {

    void setCouponPolicy(CouponPolicy couponPolicy);

    CouponPolicy getCouponPolicy(String couponPolicyId);

    void setCouponPolicyQuantity(CouponPolicy couponPolicy);

    Long getCouponPolicyQuantity(String couponPolicyId);

    Map<String, Long> getAllCouponPolicyQuantities();

    Long incrementAndGetCouponPolicyQuantity(String couponPolicyId);

    Long decrementAndGetCouponPolicyQuantity(String couponPolicyId);

    void setCouponState(CouponDTO.Response coupon);

    CouponDTO.Response getCouponState(String couponId);
}

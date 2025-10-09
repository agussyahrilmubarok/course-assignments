package com.example.coupon.service.v3;

import com.example.coupon.domain.Coupon;
import com.example.coupon.model.CouponDTO;

import java.util.List;

public interface CouponService {

    List<CouponDTO.Response> findCoupons(CouponDTO.ListRequest request);

    CouponDTO.Response findCoupon(String couponId);

    Coupon issueCoupon(CouponDTO.IssueRequest request);

    void requestIssueCoupon(CouponDTO.IssueRequest request);

    void processIssueCoupon(CouponDTO.IssueMessage message);

    Coupon useCoupon(String couponId, String orderId);

    Coupon cancelCoupon(String couponId);
}

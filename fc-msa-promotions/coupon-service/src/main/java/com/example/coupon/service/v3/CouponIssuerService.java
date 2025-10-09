package com.example.coupon.service.v3;

import com.example.coupon.domain.Coupon;
import com.example.coupon.model.CouponDTO;

public interface CouponIssuerService {

    Coupon issueCoupon(CouponDTO.IssueRequest request);
}

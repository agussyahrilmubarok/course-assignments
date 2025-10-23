package com.example.coupon.service.v2;

import com.example.coupon.domain.Coupon;
import com.example.coupon.model.CouponDTO;

public interface CouponIssuerService {

    Coupon issueCoupon(CouponDTO.IssueRequest request);
}

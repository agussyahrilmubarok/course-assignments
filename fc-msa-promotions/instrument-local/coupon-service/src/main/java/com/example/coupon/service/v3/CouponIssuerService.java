package com.example.coupon.service.v3;

import com.example.coupon.model.CouponDTO;

public interface CouponIssuerService {

    void issueCoupon(CouponDTO.IssueRequest request);

    void processIssueCoupon(CouponDTO.IssueMessage message);
}

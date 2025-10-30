package com.example.coupon.service.v3;

import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.model.CouponPolicyDTO;

import java.util.List;

public interface CouponPolicyService {

    List<CouponPolicy> findAll();

    CouponPolicy findById(String id);

    CouponPolicy create(CouponPolicyDTO.CreateRequest request);
}

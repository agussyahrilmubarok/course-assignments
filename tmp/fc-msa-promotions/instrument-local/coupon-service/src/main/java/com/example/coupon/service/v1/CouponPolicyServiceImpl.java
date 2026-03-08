package com.example.coupon.service.v1;

import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponPolicyNotFoundException;
import com.example.coupon.model.CouponPolicyDTO;
import com.example.coupon.repos.CouponPolicyRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service("CouponPolicyServiceImplV1")
@Slf4j
@RequiredArgsConstructor
public class CouponPolicyServiceImpl implements CouponPolicyService {

    private final CouponPolicyRepository couponPolicyRepository;

    /**
     * Retrieves all coupon policies.
     *
     * @return list of CouponPolicy entities
     */
    @Override
    @Transactional(readOnly = true)
    public List<CouponPolicy> findAll() {
        return couponPolicyRepository.findAll();
    }

    /**
     * Retrieves a coupon policy by ID.
     *
     * @param id coupon policy ID
     * @return CouponPolicy entity
     * @throws CouponPolicyNotFoundException if not found
     */
    @Override
    @Transactional(readOnly = true)
    public CouponPolicy findById(String id) {
        return couponPolicyRepository.findById(id)
                .orElseThrow(() -> new CouponPolicyNotFoundException("Coupon policy not found."));
    }

    /**
     * Creates and saves a new CouponPolicy from the given DTO.
     *
     * @param request DTO containing coupon policy data
     * @return saved CouponPolicy entity
     */
    @Override
    @Transactional
    public CouponPolicy create(CouponPolicyDTO.CreateRequest request) {
        CouponPolicy couponPolicy = request.toEntity();
        return couponPolicyRepository.save(couponPolicy);
    }
}

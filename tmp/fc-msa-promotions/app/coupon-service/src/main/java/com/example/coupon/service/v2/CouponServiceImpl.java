package com.example.coupon.service.v2;

import com.example.coupon.domain.Coupon;
import com.example.coupon.exception.CouponIssueException;
import com.example.coupon.exception.CouponNotFoundException;
import com.example.coupon.model.CouponDTO;
import com.example.coupon.repos.CouponRepository;
import com.example.coupon.utils.UserIdInterceptor;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Service("CouponServiceImplV2")
@Slf4j
@RequiredArgsConstructor
public class CouponServiceImpl implements CouponService {

    private final CouponRepository couponRepository;
    private final CouponIssuerService couponIssuerService;
    private final CouponRedisService couponRedisService;

    /**
     * Retrieve paginated coupons for the current user filtered by status.
     *
     * @param request filter and pagination info
     * @return list of coupons mapped to DTO responses
     */
    @Override
    @Transactional(readOnly = true)
    public List<CouponDTO.Response> findCoupons(CouponDTO.ListRequest request) {
        Pageable pageable = PageRequest.of(
                request.getPage() != null ? request.getPage() : 0,
                request.getSize() != null ? request.getSize() : 10
        );

        String currentUserId = UserIdInterceptor.getCurrentUserId();
        log.info("Fetching coupons for user: {}", currentUserId);

        List<Coupon> coupons = couponRepository
                .findByUserIdAndStatusOrderByCreatedAtDesc(currentUserId, request.getStatus(), pageable)
                .stream()
                .toList();

        log.info("Retrieved {} coupons for user {}", coupons.size(), currentUserId);

        return coupons.stream()
                .map(CouponDTO.Response::from)
                .collect(Collectors.toList());
    }

    /**
     * Retrieve a single coupon for current user by ID.
     *
     * @param couponId coupon identifier
     * @return coupon mapped to response DTO
     * @throws CouponNotFoundException if coupon not found or not owned by current user
     */
    @Override
    @Transactional(readOnly = true)
    public CouponDTO.Response findCoupon(String couponId) {
        String currentUserId = UserIdInterceptor.getCurrentUserId();
        log.info("Fetching coupon ID: {} for user {}", couponId, currentUserId);

        return couponRepository.findByIdAndUserId(couponId, currentUserId)
                .map(CouponDTO.Response::from)
                .orElseThrow(() -> {
                    log.warn("Coupon not found or unauthorized access attempt. Coupon ID: {}", couponId);
                    return new CouponNotFoundException("Coupon not found or access denied.");
                });
    }

    /**
     * Issue a coupon based on the provided coupon policy.
     * Delegates issuance to CouponIssuerService which handles concurrency and quota logic.
     *
     * @param request DTO with coupon policy ID
     * @return issued Coupon entity
     * @throws CouponIssueException if policy invalid, expired, or quota exceeded
     */
    @Override
    @Transactional
    public Coupon issueCoupon(CouponDTO.IssueRequest request) {
        log.info("Request to issue coupon for policy ID: {}", request.getCouponPolicyId());
        return couponIssuerService.issueCoupon(request);
    }

    /**
     * Mark a coupon as used for a specific order by current user.
     *
     * @param couponId coupon identifier
     * @param orderId  order identifier
     * @return updated coupon entity
     * @throws CouponNotFoundException if coupon not found or not owned by user
     */
    @Override
    @Transactional
    public Coupon useCoupon(String couponId, String orderId) {
        String currentUserId = UserIdInterceptor.getCurrentUserId();
        log.info("User {} attempting to use coupon ID: {}", currentUserId, couponId);

        Coupon coupon = couponRepository.findByIdAndUserId(couponId, currentUserId)
                .orElseThrow(() -> {
                    log.warn("Coupon not found or unauthorized use attempt. Coupon ID: {}", couponId);
                    return new CouponNotFoundException("Coupon not found or no access permission.");
                });

        coupon.use(orderId);
        Coupon updatedCoupon = couponRepository.save(coupon);
        couponRedisService.setCouponState(CouponDTO.Response.from(updatedCoupon));

        log.info("Coupon ID: {} used successfully for order ID: {}", couponId, orderId);
        return updatedCoupon;
    }

    /**
     * Cancel a previously used coupon by current user.
     *
     * @param couponId coupon identifier
     * @return updated coupon entity
     * @throws CouponNotFoundException if coupon not found or not owned by user
     */
    @Override
    @Transactional
    public Coupon cancelCoupon(String couponId) {
        String currentUserId = UserIdInterceptor.getCurrentUserId();
        log.info("User {} attempting to cancel coupon ID: {}", currentUserId, couponId);

        Coupon coupon = couponRepository.findByIdAndUserId(couponId, currentUserId)
                .orElseThrow(() -> {
                    log.warn("Coupon not found or unauthorized cancel attempt. Coupon ID: {}", couponId);
                    return new CouponNotFoundException("Coupon not found or no access permission.");
                });

        coupon.cancel();
        Coupon updatedCoupon = couponRepository.save(coupon);
        couponRedisService.incrementAndGetCouponPolicyQuantity(updatedCoupon.getCouponPolicy().getId());
        couponRedisService.setCouponState(CouponDTO.Response.from(updatedCoupon));

        log.info("Coupon ID: {} canceled successfully", couponId);
        return updatedCoupon;
    }
}
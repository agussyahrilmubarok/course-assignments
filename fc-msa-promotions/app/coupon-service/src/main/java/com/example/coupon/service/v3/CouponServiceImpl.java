package com.example.coupon.service.v3;

import com.example.coupon.domain.Coupon;
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

@Service("CouponServiceImplV3")
@Slf4j
@RequiredArgsConstructor
public class CouponServiceImpl implements CouponService {

    private final CouponRepository couponRepository;
    private final CouponIssuerService couponIssuerService;
    private final CouponRedisService couponRedisService;

    /**
     * Retrieves a paginated list of coupons for the current user filtered by status.
     *
     * @param request contains pagination and filter parameters
     * @return list of coupon DTOs belonging to the current user
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
     * Retrieves a specific coupon by ID for the current user.
     * Ensures that the coupon belongs to the user and prevents unauthorized access.
     *
     * @param couponId ID of the coupon to retrieve
     * @return coupon DTO if found
     * @throws CouponNotFoundException if the coupon does not exist or belongs to another user
     */
    @Override
    @Transactional(readOnly = true)
    public CouponDTO.Response findCoupon(String couponId) {
        String currentUserId = UserIdInterceptor.getCurrentUserId();
        log.info("Fetching coupon ID: {} for user {}", couponId, currentUserId);

        return couponRepository.findByIdAndUserId(couponId, currentUserId)
                .map(CouponDTO.Response::from)
                .orElseThrow(() -> {
                    log.error("Coupon not found or unauthorized access. Coupon ID: {}", couponId);
                    return new CouponNotFoundException("Coupon not found or no access permission.");
                });
    }

    /**
     * Issues a coupon asynchronously by sending a request to the Kafka-backed issuer service.
     *
     * @param request contains the coupon policy ID to issue a coupon from
     */
    @Override
    public void requestIssueCoupon(CouponDTO.IssueRequest request) {
        couponIssuerService.issueCoupon(request);
        log.info("Coupon issue request processed for policyId={}, userId={}",
                request.getCouponPolicyId(), UserIdInterceptor.getCurrentUserId());
    }

    /**
     * Processes the actual coupon issuance, typically called by a Kafka consumer.
     * Delegates the task to CouponIssuerService.
     *
     * @param message contains the coupon policy ID and user ID
     */
    @Override
    public void processIssueCoupon(CouponDTO.IssueMessage message) {
        log.info("Processing coupon issue message for user={}, policy={}", message.getUserId(), message.getCouponPolicyId());
        couponIssuerService.processIssueCoupon(message);
    }

    /**
     * Marks a coupon as used by the current user for a specific order.
     * Saves the updated coupon state and caches it.
     *
     * @param couponId the ID of the coupon to use
     * @param orderId  the order ID associated with the coupon usage
     * @return the updated Coupon entity
     * @throws CouponNotFoundException if the coupon does not exist or does not belong to the user
     */
    @Override
    @Transactional
    public Coupon useCoupon(String couponId, String orderId) {
        String currentUserId = UserIdInterceptor.getCurrentUserId();
        log.info("User {} is attempting to use coupon ID: {}", currentUserId, couponId);

        Coupon coupon = couponRepository.findByIdAndUserId(couponId, currentUserId)
                .orElseThrow(() -> {
                    log.error("Coupon not found or unauthorized access. Coupon ID: {}", couponId);
                    return new CouponNotFoundException("Coupon not found or no access permission.");
                });

        coupon.use(orderId);
        Coupon usedCoupon = couponRepository.save(coupon);
        couponRedisService.setCouponState(CouponDTO.Response.from(usedCoupon));
        log.info("Coupon ID: {} used successfully for order ID: {}", couponId, orderId);

        return usedCoupon;
    }

    /**
     * Cancels a coupon usage for the current user, marking it as unused again.
     * Also rolls back the available quota in Redis.
     *
     * @param couponId the ID of the coupon to cancel
     * @return the updated Coupon entity after cancellation
     * @throws CouponNotFoundException if the coupon does not exist or does not belong to the user
     */
    @Override
    @Transactional
    public Coupon cancelCoupon(String couponId) {
        String currentUserId = UserIdInterceptor.getCurrentUserId();
        log.info("User {} is attempting to cancel coupon ID: {}", currentUserId, couponId);

        Coupon coupon = couponRepository.findByIdAndUserId(couponId, currentUserId)
                .orElseThrow(() -> {
                    log.error("Coupon not found or unauthorized access. Coupon ID: {}", couponId);
                    return new CouponNotFoundException("Coupon not found or no access permission.");
                });

        coupon.cancel();
        Coupon canceledCoupon = couponRepository.save(coupon);

        // Return the coupon quota in Redis
        couponRedisService.incrementAndGetCouponPolicyQuantity(canceledCoupon.getCouponPolicy().getId());

        // Update Redis cache with latest state
        couponRedisService.setCouponState(CouponDTO.Response.from(canceledCoupon));
        log.info("Coupon ID: {} canceled successfully", couponId);

        return canceledCoupon;
    }
}

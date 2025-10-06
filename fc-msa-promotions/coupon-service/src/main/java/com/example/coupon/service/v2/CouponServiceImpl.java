package com.example.coupon.service.v2;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponIssueException;
import com.example.coupon.exception.CouponNotFoundException;
import com.example.coupon.model.CouponDTO;
import com.example.coupon.repos.CouponPolicyRepository;
import com.example.coupon.repos.CouponRepository;
import com.example.coupon.utils.UserIdInterceptor;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service("CouponServiceImplV2")
@Slf4j
@RequiredArgsConstructor
public class CouponServiceImpl implements CouponService {

    private final CouponRepository couponRepository;
    private final CouponPolicyRepository couponPolicyRepository;
    private final CouponIssuerService couponIssuerService;

    /**
     * Retrieves a paginated list of coupons for the current user based on status.
     *
     * @param request filter and pagination parameters
     * @return list of Coupon entities
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
     * Retrieves a specific coupon for the current user by its ID.
     *
     * <p>This method ensures that the coupon belongs to the current user
     * (based on {@link UserIdInterceptor}) and prevents unauthorized access
     * to other users' coupons.</p>
     *
     * @param couponId the ID of the coupon to retrieve
     * @return the coupon mapped to {@link CouponDTO.Response}
     * @throws CouponNotFoundException if coupon is not found or does not belong to the current user
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
     * Issues a coupon for the specified policy if within the valid time frame and quantity limits.
     *
     * <p><strong>Technical Limitations Based on Current Implementation:</strong></p>
     * <ul>
     *   <li><strong>1. Race Condition Risk:</strong><br>
     *   The number of issued coupons is checked using <code>countByCouponPolicyId</code> before issuing.
     *   However, there's no lock or atomic control between this check and saving the coupon.
     *   Under concurrent access, multiple transactions can pass the check and exceed <code>totalQuantity</code>.</li>
     *
     *   <li><strong>2. Performance Concern:</strong><br>
     *   The <code>countByCouponPolicyId</code> query runs for every request.
     *   As issued coupons grow, this count query may become a performance bottleneck.</li>
     *
     *   <li><strong>3. No Locking or Isolation Enforcement:</strong><br>
     *   The method does not use pessimistic or optimistic locking,
     *   making it vulnerable in high-concurrency environments where multiple coupons are issued simultaneously.</li>
     *
     *   <li><strong>4. Quantity Inaccuracy in Distributed Systems:</strong><br>
     *   Without a distributed locking mechanism or atomic counter,
     *   ensuring the total issuance stays within the limit is difficult when running across multiple instances.</li>
     * </ul>
     *
     * @param request DTO containing coupon policy ID
     * @return issued Coupon entity
     * @throws CouponIssueException if policy is invalid, expired, or quota exceeded
     */
    @Override
    @Transactional
    public Coupon issueCoupon(CouponDTO.IssueRequest request) {
        log.info("Issuing coupon for policy ID: {}", request.getCouponPolicyId());
        return couponIssuerService.issueCoupon(request);
    }

    /**
     * Marks a coupon as used for a specific order by the current user.
     *
     * @param couponId the ID of the coupon to use
     * @param orderId  the ID of the order
     * @return updated Coupon entity
     * @throws CouponNotFoundException if coupon not found or doesn't belong to user
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
        log.info("Coupon ID: {} used successfully for order ID: {}", couponId, orderId);

        return couponRepository.save(coupon);
    }

    /**
     * Cancels a previously used coupon by the current user.
     *
     * @param couponId the ID of the coupon to cancel
     * @return updated Coupon entity
     * @throws CouponNotFoundException if coupon not found or doesn't belong to user
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
        log.info("Coupon ID: {} canceled successfully", couponId);

        return couponRepository.save(coupon);
    }

    private String generateCouponCode() {
        return UUID.randomUUID().toString().substring(0, 8);
    }
}

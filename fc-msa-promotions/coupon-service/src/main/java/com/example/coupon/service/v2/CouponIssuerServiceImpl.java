package com.example.coupon.service.v2;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponIssueException;
import com.example.coupon.model.CouponDTO;
import com.example.coupon.repos.CouponPolicyRepository;
import com.example.coupon.repos.CouponRepository;
import com.example.coupon.utils.UserIdInterceptor;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

@Service("CouponIssuerServiceImplV2")
@Slf4j
@RequiredArgsConstructor
public class CouponIssuerServiceImpl implements CouponIssuerService {

    private static final String COUPON_LOCK_KEY = "coupon:lock:";
    private static final long LOCK_WAIT_TIME = 3;
    private static final long LOCK_LEASE_TIME = 5;

    private final RedissonClient redissonClient;
    private final CouponRepository couponRepository;
    private final CouponPolicyRepository couponPolicyRepository;
    private final CouponRedisService couponRedisService;

    /**
     * Issues a coupon safely with distributed lock and quota checks.
     *
     * @param request Coupon issue request DTO containing policy ID
     * @return issued Coupon entity
     * @throws CouponIssueException when locking, quota, or policy validations fail
     */
    @Override
    public Coupon issueCoupon(CouponDTO.IssueRequest request) {
        final String couponPolicyId = request.getCouponPolicyId();
        final String lockKey = COUPON_LOCK_KEY + couponPolicyId;
        final RLock lock = redissonClient.getLock(lockKey);

        // Retrieve current user once
        final String userId = UserIdInterceptor.getCurrentUserId();

        try {
            boolean isLocked = lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS);
            if (!isLocked) {
                log.warn("Failed to acquire lock for coupon policy: {}", couponPolicyId);
                throw new CouponIssueException("Too many coupon requests. Please try again later.");
            }

            final CouponPolicy couponPolicy = couponPolicyRepository.findById(couponPolicyId)
                    .orElseThrow(() -> {
                        log.warn("Coupon policy not found: {}", couponPolicyId);
                        return new CouponIssueException("Coupon policy does not exist.");
                    });

            final LocalDateTime now = LocalDateTime.now();
            final LocalDateTime startTime = couponPolicy.getStartTime();
            final LocalDateTime endTime = couponPolicy.getEndTime();

            if ((startTime != null && now.isBefore(startTime)) || (endTime != null && now.isAfter(endTime))) {
                log.warn("Coupon issuance outside valid period for policy: {}", couponPolicyId);
                throw new IllegalStateException("Coupon issuance is not within the valid period.");
            }

            boolean alreadyIssued = couponRepository.existsByUserIdAndCouponPolicyId(userId, couponPolicyId);
            if (alreadyIssued) {
                log.warn("User {} has already received a coupon for policy {}", userId, couponPolicyId);
                throw new CouponIssueException("You have already received this coupon.");
            }

            final Long remainingQuantity = couponRedisService.decrementAndGetCouponPolicyQuantity(couponPolicyId);
            if (remainingQuantity < 0) {
                couponRedisService.incrementAndGetCouponPolicyQuantity(couponPolicyId);
                log.info("Coupons exhausted for policy {}. Current quantity: {}", couponPolicyId, remainingQuantity);
                throw new CouponIssueException("All coupons have been issued.");
            }

            Coupon coupon = new Coupon();
            coupon.setId(UUID.randomUUID().toString());
            coupon.setCouponPolicy(couponPolicy);
            coupon.setCode(generateCouponCode());
            coupon.setStatus(Coupon.Status.AVAILABLE);
            coupon.setUserId(userId);

            couponRepository.save(coupon);
            couponRedisService.setCouponState(CouponDTO.Response.from(coupon));

            log.info("Coupon issued successfully. CouponId: {}, Remaining: {} for policy: {}", coupon.getId(), remainingQuantity, couponPolicyId);

            return coupon;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            log.error("Coupon issuance interrupted for policy: {}", couponPolicyId, e);
            throw new CouponIssueException("Coupon issuance interrupted.");
        } finally {
            if (lock.isHeldByCurrentThread()) {
                lock.unlock();
                log.debug("Lock released for coupon policy: {}", couponPolicyId);
            }
        }
    }

    private String generateCouponCode() {
        return UUID.randomUUID().toString().substring(0, 8);
    }
}
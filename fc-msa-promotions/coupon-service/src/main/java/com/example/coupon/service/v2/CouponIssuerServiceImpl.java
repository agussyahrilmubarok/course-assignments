package com.example.coupon.service.v2;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponIssueException;
import com.example.coupon.model.CouponDTO;
import com.example.coupon.repos.CouponRepository;
import com.example.coupon.utils.UserIdInterceptor;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RAtomicLong;
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

    private static final String COUPON_QUANTITY_KEY = "coupon:quantity:";
    private static final String COUPON_LOCK_KEY = "coupon:lock:";
    private static final long LOCK_WAIT_TIME = 3;
    private static final long LOCK_LEASE_TIME = 5;

    private final RedissonClient redissonClient;
    private final CouponRepository couponRepository;
    private final CouponPolicyService couponPolicyService;
    private final CouponRedisService couponRedisService;

    /**
     * Issues a coupon for the specified policy, ensuring no over-issuance using distributed locking.
     *
     * @param request DTO containing coupon policy ID
     * @return the Coupon used for issuance
     * @throws CouponIssueException if lock not acquired, expired, or quota exceeded
     */
    @Override
    public Coupon issueCoupon(CouponDTO.IssueRequest request) {
        String quantityKey = COUPON_QUANTITY_KEY + request.getCouponPolicyId();
        String lockKey = COUPON_LOCK_KEY + request.getCouponPolicyId();
        RLock lock = redissonClient.getLock(lockKey);

        try {
            boolean isLocked = lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS);
            if (!isLocked) {
                log.warn("Lock acquisition failed for coupon policy: {}", request.getCouponPolicyId());
                throw new CouponIssueException("Too many coupon requests. Please try again later.");
            }

            CouponPolicy couponPolicy = couponPolicyService.findById(request.getCouponPolicyId());

            LocalDateTime now = LocalDateTime.now();
            if (now.isBefore(couponPolicy.getStartTime()) || now.isAfter(couponPolicy.getEndTime())) {
                log.warn("Coupon issuance attempt outside valid period: {}", request.getCouponPolicyId());
                throw new IllegalStateException("It is not within the coupon issuance period.");
            }

            RAtomicLong atomicQuantity = redissonClient.getAtomicLong(quantityKey);
            long remainingQuantity = atomicQuantity.decrementAndGet();

            if (remainingQuantity < 0) {
                atomicQuantity.incrementAndGet();
                log.info("Coupon exhausted for policy: {}", request.getCouponPolicyId());
                throw new CouponIssueException("All coupons have been issued.");
            }

            Coupon coupon = new Coupon();
            coupon.setId(UUID.randomUUID().toString());
            coupon.setCouponPolicy(couponPolicy);
            coupon.setCode(generateCouponCode());
            coupon.setStatus(Coupon.Status.AVAILABLE);
            coupon.setUserId(UserIdInterceptor.getCurrentUserId());
            couponRepository.save(coupon);
            log.info("Save New Coupon. Coupon id: {}", coupon.getId());

            couponRedisService.setCouponState(CouponDTO.Response.from(coupon));

            log.info("Coupon issued. Remaining: {} for policy: {}", remainingQuantity, request.getCouponPolicyId());
            return coupon;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            log.error("Coupon issue interrupted for policy: {}", request.getCouponPolicyId(), e);
            throw new CouponIssueException("An error occurred while issuing the coupon.");
        } finally {
            if (lock.isHeldByCurrentThread()) {
                lock.unlock();
                log.debug("Lock released for policy: {}", request.getCouponPolicyId());
            }
        }
    }

    private String generateCouponCode() {
        return UUID.randomUUID().toString().substring(0, 8);
    }
}

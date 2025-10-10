package com.example.coupon.service.v3;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponIssueException;
import com.example.coupon.model.CouponDTO;
import com.example.coupon.repos.CouponPolicyRepository;
import com.example.coupon.repos.CouponRepository;
import com.example.coupon.service.v3.component.KafkaProducer;
import com.example.coupon.utils.UserIdInterceptor;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

@Service("CouponIssuerServiceImplV3")
@Slf4j
@RequiredArgsConstructor
public class CouponIssuerServiceImpl implements CouponIssuerService {

    // Redis key prefixes
    private static final String COUPON_LOCK_KEY = "coupon:lock:";

    // Lock configuration (seconds)
    private static final long LOCK_WAIT_TIME = 3;
    private static final long LOCK_LEASE_TIME = 5;

    private final RedissonClient redissonClient;
    private final CouponRepository couponRepository;
    private final CouponPolicyRepository couponPolicyRepository;
    private final CouponRedisService couponRedisService;
    private final KafkaProducer kafkaProducer;

    /**
     * Issues a coupon request asynchronously by first validating conditions and then sending a Kafka message.
     * Uses a distributed lock to ensure no over-issuance happens.
     *
     * @param request The coupon issue request DTO containing the policy ID.
     * @throws CouponIssueException if lock acquisition fails, the period is invalid, or quota is exceeded.
     */
    @Override
    @Transactional
    public void issueCoupon(CouponDTO.IssueRequest request) {
        String lockKey = COUPON_LOCK_KEY + request.getCouponPolicyId(); // Redis lock key
        RLock lock = redissonClient.getLock(lockKey);

        try {
            // Try to acquire lock to prevent race conditions in high-concurrency environment
            boolean isLocked = lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS);
            if (!isLocked) {
                log.warn("Failed to acquire lock for coupon policy: {}", request.getCouponPolicyId());
                throw new CouponIssueException("Too many coupon requests. Please try again later.");
            }

            // Fetch coupon policy; if not found, throw error
            CouponPolicy couponPolicy = couponPolicyRepository.findById(request.getCouponPolicyId())
                    .orElseThrow(() -> {
                        log.warn("Coupon policy not found: {}", request.getCouponPolicyId());
                        return new CouponIssueException("Coupon policy does not exist.");
                    });

            // Validate current time against policy period
            LocalDateTime now = LocalDateTime.now();
            if (now.isBefore(couponPolicy.getStartTime()) || now.isAfter(couponPolicy.getEndTime())) {
                log.warn("Coupon issuance outside valid period for policy: {}", request.getCouponPolicyId());
                throw new IllegalStateException("It is not within the coupon issuance period.");
            }

            // Get current user
            String userId = UserIdInterceptor.getCurrentUserId();

            // Check if user already has the coupon
            boolean alreadyIssued = couponRepository.existsByUserIdAndCouponPolicyId(userId, request.getCouponPolicyId());
            if (alreadyIssued) {
                log.warn("User {} already has a coupon for policy {}", userId, request.getCouponPolicyId());
                throw new CouponIssueException("You have already received this coupon.");
            }

            // Atomically decrement coupon quantity in Redis
            Long remainingQuantity = couponRedisService.decrementAndGetCouponPolicyQuantity(request.getCouponPolicyId());
            if (remainingQuantity < 0) {
                // Roll back Redis counter if exhausted
                Long rolledBack = couponRedisService.incrementAndGetCouponPolicyQuantity(request.getCouponPolicyId());
                log.info("Coupon quota exhausted for policy: {}, rolled back to: {}", request.getCouponPolicyId(), rolledBack);
                throw new CouponIssueException("All coupons have been issued.");
            }

            // Send coupon issue request to Kafka for asynchronous processing
            kafkaProducer.sendCouponIssueRequest(CouponDTO.IssueMessage.builder()
                    .couponPolicyId(request.getCouponPolicyId())
                    .userId(userId)
                    .build());

            log.info("Coupon issue request sent to Kafka for policyId: {}, userId: {}", request.getCouponPolicyId(), userId);

        } catch (InterruptedException e) {
            // Restore interrupted state and log
            Thread.currentThread().interrupt();
            log.error("Coupon issuance interrupted for policy: {}", request.getCouponPolicyId(), e);
            throw new CouponIssueException("An error occurred while issuing the coupon.");
        } finally {
            // Only unlock if the current thread holds the lock
            if (lock.isHeldByCurrentThread()) {
                lock.unlock();
                log.debug("Lock released for policy: {}", request.getCouponPolicyId());
            }
        }
    }

    /**
     * Process the actual coupon issuance synchronously (typically called from Kafka consumer).
     * This creates the coupon and persists it.
     *
     * @param message The message received from Kafka containing policy ID and user ID.
     */
    @Override
    @Transactional
    public void processIssueCoupon(CouponDTO.IssueMessage message) {
        try {
            // Find the coupon policy
            CouponPolicy couponPolicy = couponPolicyRepository.findById(message.getCouponPolicyId())
                    .orElseThrow(() -> {
                        // If not found, return coupon quantity in Redis
                        couponRedisService.incrementAndGetCouponPolicyQuantity(message.getCouponPolicyId());
                        log.warn("Coupon policy not found while processing message: {}", message.getCouponPolicyId());
                        return new CouponIssueException("Coupon policy does not exist.");
                    });

            // Create new coupon
            Coupon coupon = new Coupon();
            coupon.setId(UUID.randomUUID().toString());
            coupon.setCouponPolicy(couponPolicy);
            coupon.setCode(generateCouponCode()); // Generate unique code
            coupon.setStatus(Coupon.Status.AVAILABLE);
            coupon.setUserId(message.getUserId());

            couponRepository.save(coupon); // Save to database
            log.info("Coupon successfully issued and saved. Coupon ID: {}", coupon.getId());

            // Cache coupon in Redis for quick lookup or further processing
            couponRedisService.setCouponState(CouponDTO.Response.from(coupon));

        } catch (Exception e) {
            log.error("Failed to process coupon issuance for policyId: {}, error: {}", message.getCouponPolicyId(), e.getMessage(), e);
            // You might consider pushing this to a dead-letter queue or retry mechanism.
        }
    }

    /**
     * Generates a random coupon code (first 8 characters of a UUID).
     *
     * @return a random coupon code
     */
    private String generateCouponCode() {
        return UUID.randomUUID().toString().substring(0, 8);
    }
}

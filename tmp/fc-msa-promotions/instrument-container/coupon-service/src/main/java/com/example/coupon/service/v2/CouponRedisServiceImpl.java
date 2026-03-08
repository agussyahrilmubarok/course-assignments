package com.example.coupon.service.v2;

import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.model.CouponDTO;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RAtomicLong;
import org.redisson.api.RBucket;
import org.redisson.api.RKeys;
import org.redisson.api.RedissonClient;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;

@Service("CouponRedisServiceImplV2")
@Slf4j
@RequiredArgsConstructor
public class CouponRedisServiceImpl implements CouponRedisService {

    private static final String COUPON_POLICY_KEY = "coupon:policy:";
    private static final String COUPON_QUANTITY_KEY = "coupon:quantity:";
    private static final String COUPON_STATE_KEY = "coupon:state:";

    private final RedissonClient redissonClient;
    private final ObjectMapper objectMapper;

    /**
     * Stores the given coupon policy in Redis cache.
     *
     * @param couponPolicy the coupon policy entity to cache
     */
    @Override
    public void setCouponPolicy(CouponPolicy couponPolicy) {
        String policyKey = COUPON_POLICY_KEY + couponPolicy.getId();
        try {
            String policyJson = objectMapper.writeValueAsString(couponPolicy);
            RBucket<String> bucket = redissonClient.getBucket(policyKey);
            bucket.set(policyJson, getTTL(couponPolicy.getStartTime(), couponPolicy.getEndTime()));
            log.debug("Coupon policy cached: {}", policyKey);
        } catch (JsonProcessingException e) {
            log.error("Failed to serialize coupon policy: {}", couponPolicy.getId(), e);
            throw new RuntimeException(e);
        }
    }

    /**
     * Retrieves a cached coupon policy from Redis.
     *
     * @param couponPolicyId the ID of the coupon policy
     * @return the deserialized CouponPolicy object, or null if not found
     */
    @Override
    public CouponPolicy getCouponPolicy(String couponPolicyId) {
        String policyKey = COUPON_POLICY_KEY + couponPolicyId;
        RBucket<String> bucket = redissonClient.getBucket(policyKey);
        String policyJson = bucket.get();
        if (policyJson == null) {
            log.warn("Coupon policy not found in Redis: {}", couponPolicyId);
            return null;
        }
        try {
            return objectMapper.readValue(policyJson, CouponPolicy.class);
        } catch (JsonProcessingException e) {
            log.error("Failed to deserialize coupon policy: {}", policyJson, e);
            throw new RuntimeException("Failed to get coupon policy", e);
        }
    }

    /**
     * Sets the available quantity of a coupon policy in Redis.
     * Redis is the source of truth for coupon quantity.
     *
     * @param couponPolicy the coupon policy entity to cache
     */
    @Override
    public void setCouponPolicyQuantity(CouponPolicy couponPolicy) {
        String quantityKey = COUPON_QUANTITY_KEY + couponPolicy.getId();
        RAtomicLong atomicQuantity = redissonClient.getAtomicLong(quantityKey);
        atomicQuantity.set(couponPolicy.getTotalQuantity());
        atomicQuantity.expire(getTTL(couponPolicy.getStartTime(), couponPolicy.getEndTime()));
        log.debug("Coupon quantity set: {} = {}", quantityKey, couponPolicy.getTotalQuantity());
    }

    /**
     * Retrieves the remaining quantity of a coupon policy from Redis.
     * Redis is the source of truth for coupon quantity.
     *
     * @param couponPolicyId the ID of the coupon policy
     * @return remaining quantity of coupons
     */
    @Override
    public Long getCouponPolicyQuantity(String couponPolicyId) {
        String quantityKey = COUPON_QUANTITY_KEY + couponPolicyId;
        RAtomicLong atomicQuantity = redissonClient.getAtomicLong(quantityKey);
        return atomicQuantity.get();
    }

    /**
     * Retrieves all coupon policy quantities from Redis.
     * Returns a Map with key = couponPolicyId and value = quantity.
     *
     * @return Map<String, Long> all coupon policy quantities stored in Redis
     */
    @Override
    public Map<String, Long> getAllCouponPolicyQuantities() {
        Map<String, Long> result = new HashMap<>();
        RKeys keys = redissonClient.getKeys();

        Iterator<String> iter = keys.getKeysByPattern(COUPON_QUANTITY_KEY + "*").iterator();
        while (iter.hasNext()) {
            String fullKey = iter.next();
            String couponPolicyId = fullKey.substring(COUPON_QUANTITY_KEY.length());
            RAtomicLong atomicQuantity = redissonClient.getAtomicLong(fullKey);
            long quantity = atomicQuantity.get();
            result.put(couponPolicyId, quantity);
        }

        log.debug("Retrieved all coupon policy quantities from Redis: {}", result);
        return result;
    }

    /**
     * Increments the available coupon quantity for a given coupon policy in Redis.
     * Should be called when a coupon is canceled or released.
     *
     * @param couponPolicyId the ID of the coupon policy
     */
    @Override
    public Long incrementAndGetCouponPolicyQuantity(String couponPolicyId) {
        String quantityKey = COUPON_QUANTITY_KEY + couponPolicyId;
        RAtomicLong atomicQuantity = redissonClient.getAtomicLong(quantityKey);
        return atomicQuantity.incrementAndGet();
    }

    /**
     * Decrements the available coupon quantity for a given coupon policy in Redis.
     * Should be called when a coupon is issued.
     *
     * @param couponPolicyId the ID of the coupon policy
     */
    @Override
    public Long decrementAndGetCouponPolicyQuantity(String couponPolicyId) {
        String quantityKey = COUPON_QUANTITY_KEY + couponPolicyId;
        RAtomicLong atomicQuantity = redissonClient.getAtomicLong(quantityKey);
        return atomicQuantity.decrementAndGet();
    }

    /**
     * Caches the current state of a coupon in Redis.
     *
     * @param coupon the coupon entity to store
     */
    @Override
    public void setCouponState(CouponDTO.Response coupon) {
        String stateKey = COUPON_STATE_KEY + coupon.getId();
        try {
            String couponJson = objectMapper.writeValueAsString(coupon);
            RBucket<String> bucket = redissonClient.getBucket(stateKey);
            bucket.set(couponJson, getTTL(coupon.getValidFrom(), coupon.getValidUntil()));
            log.debug("Coupon state cached: {}", stateKey);
        } catch (JsonProcessingException e) {
            log.error("Failed to save coupon state: {}", coupon.getId(), e);
            throw new RuntimeException("Failed to save coupon", e);
        }
    }

    /**
     * Retrieves the current state of a coupon from Redis.
     *
     * @param couponId the ID of the coupon
     * @return the deserialized Coupon object, or null if not found
     */
    @Override
    public CouponDTO.Response getCouponState(String couponId) {
        String stateKey = COUPON_STATE_KEY + couponId;
        RBucket<String> bucket = redissonClient.getBucket(stateKey);
        String couponJson = bucket.get();
        if (couponJson == null) {
            return null;
        }
        try {
            return objectMapper.readValue(couponJson, CouponDTO.Response.class);
        } catch (JsonProcessingException e) {
            throw new RuntimeException("Failed to get coupon", e);
        }
    }

    private Duration getTTL(LocalDateTime start, LocalDateTime end) {
        return Duration.between(start, end);
    }
}

package com.example.coupon.service.v3;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.model.CouponDTO;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.redisson.api.*;

import java.time.Duration;
import java.time.LocalDateTime;
import java.util.Map;
import java.util.Set;

import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CouponRedisServiceImplTest {

    private static final String TEST_POLICY_ID = "COUPON_POLICY_1";
    private static final String TEST_USER_ID = "COUPON_USER_1";
    private static final String TEST_COUPON_ID = "COUPON_1";
    private static final Long TEST_QUANTITY = 100L;

    @InjectMocks
    private CouponRedisServiceImpl couponRedisService;

    @Mock
    private RedissonClient redissonClient;

    @Mock
    private ObjectMapper objectMapper;

    @Mock
    private RBucket<Object> mockObjectBucket;

    @Mock
    private RAtomicLong mockAtomicLong;

    @Mock
    private RLock mockLock;

    @Mock
    private RKeys mockKeys;

    private CouponPolicy couponPolicy;
    private Coupon coupon;

    @BeforeEach
    void setUp() {
        couponPolicy = new CouponPolicy();
        couponPolicy.setId(TEST_POLICY_ID);
        couponPolicy.setName("Test Coupon");
        couponPolicy.setDiscountType(CouponPolicy.DiscountType.FIXED_AMOUNT);
        couponPolicy.setDiscountValue(1000);
        couponPolicy.setMinimumOrderAmount(10000);
        couponPolicy.setMaximumDiscountAmount(1000);
        couponPolicy.setTotalQuantity(100);
        couponPolicy.setStartTime(LocalDateTime.now().minusDays(1));
        couponPolicy.setEndTime(LocalDateTime.now().plusDays(1));

        coupon = new Coupon();
        coupon.setId(TEST_COUPON_ID);
        coupon.setUserId(TEST_USER_ID);
        coupon.setCouponPolicy(couponPolicy);
        coupon.setCode("TEST123");
    }

    @Test
    void givenValidCouponPolicy_whenSetCouponPolicy_thenCacheIsSaved() throws JsonProcessingException {
        String expectedJson = "{\"id\":\"COUPON_POLICY_1\",\"name\":\"Test Coupon\"}";
        Duration expectedTTL = Duration.ofHours(48);

        when(objectMapper.writeValueAsString(couponPolicy)).thenReturn(expectedJson);
        when(redissonClient.getBucket("coupon:policy:" + TEST_POLICY_ID)).thenReturn(mockObjectBucket);

        couponRedisService.setCouponPolicy(couponPolicy);

        verify(mockObjectBucket).set(eq(expectedJson), argThat(ttl -> ttl.compareTo(expectedTTL) >= 0));
    }

    @Test
    void givenInvalidCouponPolicy_whenSetCouponPolicy_thenCacheIsNotSaved() throws JsonProcessingException {
        when(objectMapper.writeValueAsString(couponPolicy)).thenThrow(JsonProcessingException.class);

        assertThrows(RuntimeException.class, () -> {
            couponRedisService.setCouponPolicy(couponPolicy);
        });
    }

    @Test
    void givenValidCouponPolicyInRedis_whenGetCouponPolicy_thenReturnCouponPolicy() throws JsonProcessingException {
        String policyJson = "{\"id\":\"COUPON_POLICY_1\",\"name\":\"Test Coupon\"}";

        when(redissonClient.getBucket("coupon:policy:" + TEST_POLICY_ID)).thenReturn(mockObjectBucket);
        when(mockObjectBucket.get()).thenReturn(policyJson);
        when(objectMapper.readValue(policyJson, CouponPolicy.class)).thenReturn(couponPolicy);

        CouponPolicy result = couponRedisService.getCouponPolicy(TEST_POLICY_ID);

        assertNotNull(result);
        assertEquals(TEST_POLICY_ID, result.getId());
        verify(objectMapper).readValue(policyJson, CouponPolicy.class);
    }

    @Test
    void givenNoCouponPolicyInRedis_whenGetCouponPolicy_thenReturnNull() {
        when(redissonClient.getBucket("coupon:policy:" + TEST_POLICY_ID)).thenReturn(mockObjectBucket);
        when(mockObjectBucket.get()).thenReturn(null);

        CouponPolicy result = couponRedisService.getCouponPolicy(TEST_POLICY_ID);

        assertNull(result);
    }

    @Test
    void givenNoCouponPolicyInRedis_whenGetCouponPolicy_thenReturnDeserialization() throws JsonProcessingException {
        String policyJson = "{\"id\":\"COUPON_POLICY_1\",\"name\":\"Test Coupon\"}";

        when(redissonClient.getBucket("coupon:policy:" + TEST_POLICY_ID)).thenReturn(mockObjectBucket);
        when(mockObjectBucket.get()).thenReturn(policyJson);
        when(objectMapper.readValue(policyJson, CouponPolicy.class)).thenThrow(JsonProcessingException.class);

        assertThrows(RuntimeException.class, () -> couponRedisService.getCouponPolicy(TEST_POLICY_ID));
    }

    @Test
    void givenCouponPolicy_whenSetCouponPolicyQuantity_thenQuantityIsSet() {
        when(redissonClient.getAtomicLong("coupon:quantity:" + TEST_POLICY_ID)).thenReturn(mockAtomicLong);

        couponRedisService.setCouponPolicyQuantity(couponPolicy);

        verify(mockAtomicLong).set(TEST_QUANTITY);
    }

    @Test
    void givenCouponPolicyQuantityInRedis_whenGetCouponPolicyQuantity_thenReturnQuantity() {
        when(redissonClient.getAtomicLong("coupon:quantity:" + TEST_POLICY_ID)).thenReturn(mockAtomicLong);
        when(mockAtomicLong.get()).thenReturn(TEST_QUANTITY);

        Long result = couponRedisService.getCouponPolicyQuantity(TEST_POLICY_ID);

        assertEquals(TEST_QUANTITY, result);
        verify(mockAtomicLong).get();
    }

    @Test
    void givenCouponQuantitiesInRedis_whenGetAllCouponPolicyQuantities_thenReturnAllQuantities() {
        String key1 = "coupon:quantity:POLICY1";
        String key2 = "coupon:quantity:POLICY2";

        when(redissonClient.getKeys()).thenReturn(mockKeys);
        when(mockKeys.getKeysByPattern("coupon:quantity:*")).thenReturn(Set.of(key1, key2));
        RAtomicLong atomicLong1 = mock(RAtomicLong.class);
        RAtomicLong atomicLong2 = mock(RAtomicLong.class);
        when(redissonClient.getAtomicLong(key1)).thenReturn(atomicLong1);
        when(redissonClient.getAtomicLong(key2)).thenReturn(atomicLong2);
        when(atomicLong1.get()).thenReturn(50L);
        when(atomicLong2.get()).thenReturn(25L);

        Map<String, Long> quantities = couponRedisService.getAllCouponPolicyQuantities();

        assertEquals(2, quantities.size());
        assertEquals(50L, quantities.get("POLICY1"));
        assertEquals(25L, quantities.get("POLICY2"));
    }

    @Test
    void givenPolicyId_whenDecrementAndGetQuantity_thenReturnDecrementedValue() {
        when(redissonClient.getAtomicLong("coupon:quantity:" + TEST_POLICY_ID)).thenReturn(mockAtomicLong);
        when(mockAtomicLong.decrementAndGet()).thenReturn(TEST_QUANTITY - 1);

        Long result = couponRedisService.decrementAndGetCouponPolicyQuantity(TEST_POLICY_ID);

        assertEquals(TEST_QUANTITY - 1, result);
        verify(mockAtomicLong).decrementAndGet();
    }

    @Test
    void givenPolicyId_whenIncrementAndGetQuantity_thenReturnIncrementedValue() {
        when(redissonClient.getAtomicLong("coupon:quantity:" + TEST_POLICY_ID)).thenReturn(mockAtomicLong);
        when(mockAtomicLong.incrementAndGet()).thenReturn(TEST_QUANTITY + 1);

        Long result = couponRedisService.incrementAndGetCouponPolicyQuantity(TEST_POLICY_ID);

        assertEquals(TEST_QUANTITY + 1, result);
        verify(mockAtomicLong).incrementAndGet();
    }

    @Test
    void givenCouponResponse_whenSetCouponState_thenCacheIsSaved() throws JsonProcessingException {
        String expectedJson = "{\"id\":\"COUPON_1\",\"code\":\"TEST123\"}";
        CouponDTO.Response newCoupon = CouponDTO.Response.from(coupon);
        Duration expectedTTL = Duration.between(newCoupon.getValidFrom(), newCoupon.getValidUntil());

        when(objectMapper.writeValueAsString(any(CouponDTO.Response.class))).thenReturn(expectedJson);
        when(redissonClient.getBucket("coupon:state:" + TEST_COUPON_ID)).thenReturn(mockObjectBucket);

        couponRedisService.setCouponState(CouponDTO.Response.from(coupon));

        verify(mockObjectBucket).set(eq(expectedJson), argThat(ttl -> ttl.compareTo(expectedTTL) >= 0));
    }

    @Test
    void givenCouponResponse_whenSetCouponState_thenCacheIsFail() throws Exception {
        when(objectMapper.writeValueAsString(any())).thenThrow(new JsonProcessingException("error") {
        });

        assertThatThrownBy(() -> couponRedisService.setCouponState(CouponDTO.Response.from(coupon)))
                .isInstanceOf(RuntimeException.class)
                .hasMessageContaining("Failed to save coupon");

        verify(objectMapper).writeValueAsString(any(CouponDTO.Response.class));
        verifyNoInteractions(mockObjectBucket);
    }

    @Test
    void givenExistingCouponStateInRedis_whenGetCouponState_thenReturnCouponResponse() throws JsonProcessingException {
        String couponJson = "{\"id\":\"COUPON_1\",\"code\":\"TEST123\"}";

        when(redissonClient.getBucket("coupon:state:" + TEST_COUPON_ID)).thenReturn(mockObjectBucket);
        when(mockObjectBucket.get()).thenReturn(couponJson);
        when(objectMapper.readValue(couponJson, CouponDTO.Response.class)).thenReturn(CouponDTO.Response.from(coupon));

        CouponDTO.Response result = couponRedisService.getCouponState(TEST_COUPON_ID);

        assert result != null;
        assert result.getCouponCode().equals("TEST123");
    }

    @Test
    void givenNoCouponStateInRedis_whenGetCouponState_thenReturnNull() {
        when(redissonClient.getBucket("coupon:state:" + TEST_COUPON_ID)).thenReturn(mockObjectBucket);
        when(mockObjectBucket.get()).thenReturn(null);

        CouponDTO.Response result = couponRedisService.getCouponState(TEST_COUPON_ID);

        assertNull(result);
    }
}
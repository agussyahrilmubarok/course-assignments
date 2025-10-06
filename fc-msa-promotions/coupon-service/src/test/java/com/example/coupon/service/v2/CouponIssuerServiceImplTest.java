package com.example.coupon.service.v2;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponIssueException;
import com.example.coupon.model.CouponDTO;
import com.example.coupon.repos.CouponRepository;
import com.example.coupon.utils.UserIdInterceptor;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockedStatic;
import org.mockito.junit.jupiter.MockitoExtension;
import org.redisson.api.RAtomicLong;
import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;

import java.time.LocalDateTime;
import java.util.concurrent.TimeUnit;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyLong;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CouponIssuerServiceImplTest {

    private static final String TEST_POLICY_ID = "COUPON_POLICY_1";
    private static final String TEST_USER_ID = "COUPON_USER_1";
    private static final String TEST_COUPON_ID = "COUPON_1";

    @InjectMocks
    private CouponIssuerServiceImpl couponIssuerService;

    @Mock
    private RedissonClient redissonClient;

    @Mock
    private CouponRepository couponRepository;

    @Mock
    private CouponRedisService couponRedisService;

    @Mock
    private CouponPolicyService couponPolicyService;

    @Mock
    private RAtomicLong mockAtomicLong;

    @Mock
    private RLock mockLock;

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
    void givenValidRequest_whenIssueCoupon_thenReturnCoupon() throws Exception {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any(TimeUnit.class))).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(redissonClient.getAtomicLong("coupon:quantity:" + TEST_POLICY_ID)).thenReturn(mockAtomicLong);
        when(mockAtomicLong.decrementAndGet()).thenReturn(99L);
        when(couponPolicyService.findById(TEST_POLICY_ID)).thenReturn(couponPolicy);
        when(couponRepository.save(any(Coupon.class))).thenAnswer(invocation -> {
            Coupon c = invocation.getArgument(0);
            c.setId(TEST_COUPON_ID);
            return c;
        });

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            Coupon issuedCoupon = couponIssuerService.issueCoupon(request);

            assertThat(issuedCoupon.getId()).isEqualTo(TEST_COUPON_ID);
            assertThat(issuedCoupon.getUserId()).isEqualTo(TEST_USER_ID);
            verify(couponRepository).save(any(Coupon.class));
            verify(mockLock).unlock();
        }
    }

    @Test
    void givenValidRequest_whenIssueCouponLockFail_thenReturnException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(false);

        assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));
    }

    @Test
    void givenInvalidPeriodRequest_whenIssueCoupon_thenReturnException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();
        couponPolicy.setStartTime(LocalDateTime.now().plusDays(1));

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyService.findById(TEST_POLICY_ID)).thenReturn(couponPolicy);

        assertThrows(IllegalStateException.class, () -> couponIssuerService.issueCoupon(request));
        verify(mockLock).unlock();
    }

    @Test
    void givenCouponQuotaExhausted_whenIssueCoupon_thenReturnException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyService.findById(TEST_POLICY_ID)).thenReturn(couponPolicy);
        when(redissonClient.getAtomicLong("coupon:quantity:" + TEST_POLICY_ID)).thenReturn(mockAtomicLong);
        when(mockAtomicLong.decrementAndGet()).thenReturn(-1L);

        assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));
        verify(mockAtomicLong).incrementAndGet(); // rollback
        verify(mockLock).unlock();
    }
}
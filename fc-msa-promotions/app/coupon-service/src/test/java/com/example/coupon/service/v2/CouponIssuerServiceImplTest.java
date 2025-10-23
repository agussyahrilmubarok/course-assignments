package com.example.coupon.service.v2;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponIssueException;
import com.example.coupon.model.CouponDTO;
import com.example.coupon.repos.CouponPolicyRepository;
import com.example.coupon.repos.CouponRepository;
import com.example.coupon.utils.UserIdInterceptor;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockedStatic;
import org.mockito.junit.jupiter.MockitoExtension;
import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;

import java.time.LocalDateTime;
import java.util.Optional;
import java.util.concurrent.TimeUnit;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.*;
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
    private CouponPolicyRepository couponPolicyRepository;
    @Mock
    private CouponRedisService couponRedisService;
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
    void testIssueCoupon_whenValidRequest_shouldReturnCoupon() throws Exception {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any(TimeUnit.class))).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));
        when(couponRepository.existsByUserIdAndCouponPolicyId(TEST_USER_ID, TEST_POLICY_ID)).thenReturn(false);
        when(couponRedisService.decrementAndGetCouponPolicyQuantity(TEST_POLICY_ID)).thenReturn(99L);
        when(couponRepository.save(any(Coupon.class))).thenAnswer(invocation -> {
            Coupon c = invocation.getArgument(0);
            c.setId(TEST_COUPON_ID);
            return c;
        });
        doNothing().when(couponRedisService).setCouponState(any());

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
    void testIssueCoupon_whenLockFails_shouldThrowCouponIssueException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(false);

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));
        }
    }

    @Test
    void testIssueCoupon_whenBeforeStartPeriod_shouldThrowIllegalStateException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();
        couponPolicy.setStartTime(LocalDateTime.now().plusDays(1)); // start time di masa depan

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            IllegalStateException exception = assertThrows(IllegalStateException.class, () -> couponIssuerService.issueCoupon(request));
            assertThat(exception.getMessage()).isEqualTo("Coupon issuance is not within the valid period.");

            verify(mockLock).unlock();
        }
    }

    @Test
    void testIssueCoupon_whenAfterEndPeriod_shouldThrowIllegalStateException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        couponPolicy.setEndTime(LocalDateTime.now().minusSeconds(1));

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            IllegalStateException exception = assertThrows(IllegalStateException.class, () ->
                    couponIssuerService.issueCoupon(request)
            );

            assertThat(exception.getMessage()).isEqualTo("Coupon issuance is not within the valid period.");
            verify(mockLock).unlock();
        }
    }

    @Test
    void testIssueCoupon_whenUserAlreadyHasCoupon_shouldThrowCouponIssueException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(couponRepository.existsByUserIdAndCouponPolicyId(TEST_USER_ID, TEST_POLICY_ID)).thenReturn(true);

            CouponIssueException exception = assertThrows(CouponIssueException.class, () ->
                    couponIssuerService.issueCoupon(request)
            );

            assertThat(exception.getMessage()).isEqualTo("You have already received this coupon.");
            verify(mockLock).unlock();
        }
    }

    @Test
    void testIssueCoupon_whenQuotaExhausted_shouldThrowCouponIssueException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(couponRepository.existsByUserIdAndCouponPolicyId(TEST_USER_ID, TEST_POLICY_ID)).thenReturn(false);
            when(couponRedisService.decrementAndGetCouponPolicyQuantity(TEST_POLICY_ID)).thenReturn(-1L);
            when(couponRedisService.incrementAndGetCouponPolicyQuantity(TEST_POLICY_ID)).thenReturn(0L);

            assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));

            verify(couponRedisService).incrementAndGetCouponPolicyQuantity(TEST_POLICY_ID);
            verify(mockLock).unlock();
        }
    }

    @Test
    void testIssueCoupon_whenInterruptedException_shouldThrowCouponIssueException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenThrow(new InterruptedException());
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));

            verify(mockLock).unlock();
        }
    }

    @Test
    void testIssueCoupon_whenLockNotHeld_shouldNotUnlockLock() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(false);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));
        when(couponRepository.existsByUserIdAndCouponPolicyId(anyString(), anyString())).thenReturn(false);
        when(couponRedisService.decrementAndGetCouponPolicyQuantity(TEST_POLICY_ID)).thenReturn(99L);
        when(couponRepository.save(any(Coupon.class))).thenAnswer(invocation -> {
            Coupon c = invocation.getArgument(0);
            c.setId(TEST_COUPON_ID);
            return c;
        });
        doNothing().when(couponRedisService).setCouponState(any());

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            Coupon issuedCoupon = couponIssuerService.issueCoupon(request);

            assertThat(issuedCoupon.getId()).isEqualTo(TEST_COUPON_ID);
            assertThat(issuedCoupon.getUserId()).isEqualTo(TEST_USER_ID);
            verify(mockLock, never()).unlock();
        }
    }
}
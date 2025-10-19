package com.example.coupon.service.v3;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponIssueException;
import com.example.coupon.model.CouponDTO;
import com.example.coupon.repos.CouponPolicyRepository;
import com.example.coupon.repos.CouponRepository;
import com.example.coupon.service.v3.component.KafkaProducer;
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
    private KafkaProducer kafkaProducer;

    @Mock
    private RLock mockLock;

    private CouponPolicy couponPolicy;

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
    }

    @Test
    void givenValidRequest_whenIssueCoupon_thenSendKafkaMessage() throws Exception {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock("coupon:lock:" + TEST_POLICY_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any(TimeUnit.class))).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));
        when(couponRepository.existsByUserIdAndCouponPolicyId(TEST_USER_ID, TEST_POLICY_ID)).thenReturn(false);
        when(couponRedisService.decrementAndGetCouponPolicyQuantity(TEST_POLICY_ID)).thenReturn(99L);

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            couponIssuerService.issueCoupon(request);

            verify(kafkaProducer).sendCouponIssueRequest(argThat(message ->
                    message.getCouponPolicyId().equals(TEST_POLICY_ID) &&
                            message.getUserId().equals(TEST_USER_ID)
            ));
            verify(mockLock).unlock();
        }
    }

    @Test
    void givenLockNotAcquired_whenIssueCoupon_thenThrowException() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .build();

        when(redissonClient.getLock(anyString())).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(false);

        assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));
    }

    @Test
    void givenPolicyNotFound_whenIssueCoupon_thenThrowException() throws Exception {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder().couponPolicyId(TEST_POLICY_ID).build();

        when(redissonClient.getLock(anyString())).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.empty());

        assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));
        verify(mockLock).unlock();
    }

    @Test
    void givenInvalidPeriod_whenIssueCoupon_thenThrowIllegalStateException() throws Exception {
        couponPolicy.setStartTime(LocalDateTime.now().plusDays(1)); // Not started yet

        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder().couponPolicyId(TEST_POLICY_ID).build();

        when(redissonClient.getLock(anyString())).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));

        assertThrows(IllegalStateException.class, () -> couponIssuerService.issueCoupon(request));
        verify(mockLock).unlock();
    }

    @Test
    void givenUserAlreadyHasCoupon_whenIssueCoupon_thenThrowException() throws Exception {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder().couponPolicyId(TEST_POLICY_ID).build();

        when(redissonClient.getLock(anyString())).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));

        try (MockedStatic<UserIdInterceptor> mocked = mockStatic(UserIdInterceptor.class)) {
            mocked.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(couponRepository.existsByUserIdAndCouponPolicyId(TEST_USER_ID, TEST_POLICY_ID)).thenReturn(true);

            CouponIssueException ex = assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));
            assertThat(ex.getMessage()).isEqualTo("You have already received this coupon.");
            verify(mockLock).unlock();
        }
    }

    @Test
    void givenQuotaExhausted_whenIssueCoupon_thenThrowExceptionAndRollbackQuota() throws Exception {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder().couponPolicyId(TEST_POLICY_ID).build();

        when(redissonClient.getLock(anyString())).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(true);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));

        try (MockedStatic<UserIdInterceptor> mocked = mockStatic(UserIdInterceptor.class)) {
            mocked.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(couponRepository.existsByUserIdAndCouponPolicyId(TEST_USER_ID, TEST_POLICY_ID)).thenReturn(false);
            when(couponRedisService.decrementAndGetCouponPolicyQuantity(TEST_POLICY_ID)).thenReturn(-1L);
            when(couponRedisService.incrementAndGetCouponPolicyQuantity(TEST_POLICY_ID)).thenReturn(0L);

            assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));

            verify(couponRedisService).incrementAndGetCouponPolicyQuantity(TEST_POLICY_ID);
            verify(mockLock).unlock();
        }
    }

    @Test
    void givenInterruptedException_whenIssueCoupon_thenThrowExceptionAndUnlock() throws InterruptedException {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder().couponPolicyId(TEST_POLICY_ID).build();

        when(redissonClient.getLock(anyString())).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenThrow(new InterruptedException());
        when(mockLock.isHeldByCurrentThread()).thenReturn(true); // unlockable

        assertThrows(CouponIssueException.class, () -> couponIssuerService.issueCoupon(request));
        verify(mockLock).unlock();
    }

    @Test
    void givenLockNotHeldByThread_whenIssueCoupon_thenDoNotUnlock() throws Exception {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder().couponPolicyId(TEST_POLICY_ID).build();

        when(redissonClient.getLock(anyString())).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockLock.isHeldByCurrentThread()).thenReturn(false);
        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));
        when(couponRepository.existsByUserIdAndCouponPolicyId(TEST_USER_ID, TEST_POLICY_ID)).thenReturn(false);
        when(couponRedisService.decrementAndGetCouponPolicyQuantity(TEST_POLICY_ID)).thenReturn(99L);

        try (MockedStatic<UserIdInterceptor> mocked = mockStatic(UserIdInterceptor.class)) {
            mocked.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            couponIssuerService.issueCoupon(request);

            verify(mockLock, never()).unlock();
            verify(kafkaProducer).sendCouponIssueRequest(any());
        }
    }

    @Test
    void givenValidMessage_whenProcessIssueCoupon_thenSaveCouponAndCacheIt() {
        CouponDTO.IssueMessage message = CouponDTO.IssueMessage.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .userId(TEST_USER_ID)
                .build();

        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));

        couponIssuerService.processIssueCoupon(message);

        verify(couponRepository, times(1)).save(any(Coupon.class));
        verify(couponRedisService, times(1)).setCouponState(any(CouponDTO.Response.class));
    }

    @Test
    void givenPolicyNotFound_whenProcessIssueCoupon_thenRollbackQuota() {
        CouponDTO.IssueMessage message = CouponDTO.IssueMessage.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .userId(TEST_USER_ID)
                .build();

        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.empty());

        couponIssuerService.processIssueCoupon(message);

        verify(couponRedisService).incrementAndGetCouponPolicyQuantity(TEST_POLICY_ID);
        verify(couponRepository, never()).save(any());
    }

    @Test
    void givenExceptionOccurs_whenProcessIssueCoupon_thenLogError() {
        CouponDTO.IssueMessage message = CouponDTO.IssueMessage.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .userId(TEST_USER_ID)
                .build();

        when(couponPolicyRepository.findById(TEST_POLICY_ID)).thenReturn(Optional.of(couponPolicy));
        doThrow(new RuntimeException("DB error")).when(couponRepository).save(any());

        couponIssuerService.processIssueCoupon(message);

        verify(couponRepository).save(any());
        verify(couponRedisService, never()).setCouponState(any());
    }
}
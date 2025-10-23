package com.example.coupon.service.v3;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponNotFoundException;
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
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;

import java.time.LocalDateTime;
import java.util.Collections;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.empty;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CouponServiceImplTest {

    private static final String TEST_USER_ID = "USER_1";
    private static final String TEST_COUPON_POLICY_ID = "COUPON_POLICY_1";
    private static final String TEST_COUPON_ID = "COUPON_1";
    private static final String TEST_ORDER_ID = "ORDER_1";

    @InjectMocks
    private CouponServiceImpl couponService;

    @Mock
    private CouponRepository couponRepository;
    @Mock
    private CouponIssuerService couponIssuerService;
    @Mock
    private CouponRedisService couponRedisService;

    private CouponPolicy couponPolicy;
    private Coupon coupon;

    @BeforeEach
    void setUp() {
        couponPolicy = new CouponPolicy();
        couponPolicy.setId(TEST_COUPON_POLICY_ID);
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
        coupon.setOrderId(TEST_ORDER_ID);
        coupon.setCode("TEST123");
    }

    @Test
    void testFindCoupons_whenValidRequest_shouldReturnCouponList() {
        List<Coupon> coupons = List.of(coupon);
        Page<Coupon> couponPage = new PageImpl<>(coupons);
        CouponDTO.ListRequest request = CouponDTO.ListRequest.builder()
                .status(Coupon.Status.AVAILABLE)
                .page(0)
                .size(10)
                .build();

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            when(couponRepository.findByUserIdAndStatusOrderByCreatedAtDesc(eq(TEST_USER_ID), any(), any(PageRequest.class)))
                    .thenReturn(couponPage);

            List<CouponDTO.Response> results = couponService.findCoupons(request);

            assertThat(results.getFirst().getId()).isEqualTo(TEST_COUPON_ID);
            assertThat(results.getFirst().getUserId()).isEqualTo(TEST_USER_ID);
        }
    }

    @Test
    void testFindCoupons_whenNoCoupons_shouldReturnEmptyList() {
        Page<Coupon> couponPage = new PageImpl<>(Collections.emptyList());
        CouponDTO.ListRequest request = CouponDTO.ListRequest.builder()
                .status(Coupon.Status.AVAILABLE)
                .page(0)
                .size(10)
                .build();

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            when(couponRepository.findByUserIdAndStatusOrderByCreatedAtDesc(eq(TEST_USER_ID), any(), any(PageRequest.class)))
                    .thenReturn(couponPage);

            List<CouponDTO.Response> results = couponService.findCoupons(request);

            assertThat(results, empty());
        }
    }

    @Test
    void testFindCoupon_whenCouponExists_shouldReturnCoupon() {
        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            when(couponRepository.findByIdAndUserId(TEST_COUPON_ID, TEST_USER_ID)).thenReturn(Optional.of(coupon));

            CouponDTO.Response result = couponService.findCoupon(TEST_COUPON_ID);

            assertThat(result.getId()).isEqualTo(TEST_COUPON_ID);
        }
    }

    @Test
    void testFindCoupon_whenCouponNotFound_shouldThrowException() {
        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            when(couponRepository.findByIdAndUserId(TEST_COUPON_ID, TEST_USER_ID)).thenReturn(Optional.empty());

            assertThrows(CouponNotFoundException.class,
                    () -> couponService.findCoupon(TEST_COUPON_ID));
        }
    }

    @Test
    void testRequestIssueCoupon_whenCalled_shouldDelegateToIssuerService() {
        CouponDTO.IssueRequest request = CouponDTO.IssueRequest.builder()
                .couponPolicyId(TEST_COUPON_POLICY_ID)
                .build();

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            couponService.requestIssueCoupon(request);

            verify(couponIssuerService).issueCoupon(eq(request));
        }
    }

    @Test
    void testProcessIssueCoupon_whenMessageReceived_shouldDelegateToIssuerService() {
        CouponDTO.IssueMessage message = CouponDTO.IssueMessage.builder()
                .couponPolicyId(TEST_COUPON_POLICY_ID)
                .userId(TEST_USER_ID)
                .build();

        couponService.processIssueCoupon(message);

        verify(couponIssuerService).processIssueCoupon(eq(message));
    }

    @Test
    void testUseCoupon_whenCouponIsValid_shouldMarkAsUsedAndSave() {
        when(couponRepository.findByIdAndUserId(TEST_COUPON_ID, TEST_USER_ID)).thenReturn(Optional.of(coupon));
        when(couponRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            Coupon usedCoupon = couponService.useCoupon(TEST_COUPON_ID, TEST_ORDER_ID);

            assertEquals(Coupon.Status.USED, usedCoupon.getStatus());
            assertEquals(TEST_ORDER_ID, usedCoupon.getOrderId());

            verify(couponRepository).save(usedCoupon);
            verify(couponRedisService).setCouponState(any(CouponDTO.Response.class));
        }
    }

    @Test
    void testUseCoupon_whenCouponNotFound_shouldThrowException() {
        when(couponRepository.findByIdAndUserId(TEST_COUPON_ID, TEST_USER_ID)).thenReturn(Optional.empty());

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            assertThrows(CouponNotFoundException.class,
                    () -> couponService.useCoupon(TEST_COUPON_ID, TEST_ORDER_ID));
        }
    }

    @Test
    void testCancelCoupon_whenUsedCoupon_shouldMarkAsCanceledAndSave() {
        coupon.setStatus(Coupon.Status.USED);
        when(couponRepository.findByIdAndUserId(TEST_COUPON_ID, TEST_USER_ID)).thenReturn(Optional.of(coupon));
        when(couponRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            Coupon canceled = couponService.cancelCoupon(TEST_COUPON_ID);

            assertEquals(Coupon.Status.CANCELED, canceled.getStatus());

            verify(couponRepository).save(canceled);
            verify(couponRedisService).incrementAndGetCouponPolicyQuantity(TEST_COUPON_POLICY_ID);
            verify(couponRedisService).setCouponState(any(CouponDTO.Response.class));
        }
    }

    @Test
    void testCancelCoupon_whenCouponNotFound_shouldThrowException() {
        when(couponRepository.findByIdAndUserId(TEST_COUPON_ID, TEST_USER_ID)).thenReturn(Optional.empty());

        try (MockedStatic<UserIdInterceptor> mockedStatic = mockStatic(UserIdInterceptor.class)) {
            mockedStatic.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            assertThrows(CouponNotFoundException.class,
                    () -> couponService.cancelCoupon(TEST_COUPON_ID));
        }
    }
}
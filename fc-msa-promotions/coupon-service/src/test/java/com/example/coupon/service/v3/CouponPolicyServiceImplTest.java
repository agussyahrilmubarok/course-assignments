package com.example.coupon.service.v3;

import com.example.coupon.domain.CouponPolicy;
import com.example.coupon.exception.CouponPolicyNotFoundException;
import com.example.coupon.model.CouponPolicyDTO;
import com.example.coupon.repos.CouponPolicyRepository;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class CouponPolicyServiceImplTest {

    @InjectMocks
    private CouponPolicyServiceImpl couponPolicyService;

    @Mock
    private CouponPolicyRepository couponPolicyRepository;

    private CouponPolicy couponPolicy;

    @Test
    void givenExistingPolicies_whenFindAll_thenReturnListOfPolicies() {
        CouponPolicy policy1 = new CouponPolicy();
        policy1.setId("1");
        CouponPolicy policy2 = new CouponPolicy();
        policy2.setId("2");

        when(couponPolicyRepository.findAll()).thenReturn(Arrays.asList(policy1, policy2));

        List<CouponPolicy> result = couponPolicyService.findAll();

        assertThat(result).hasSize(2);
        assertThat(result.get(0).getId()).isEqualTo("1");
        assertThat(result.get(1).getId()).isEqualTo("2");
    }

    @Test
    void givenExistingId_whenFindById_thenReturnCouponPolicy() {
        CouponPolicy policy = new CouponPolicy();
        policy.setId("123");

        when(couponPolicyRepository.findById("123")).thenReturn(Optional.of(policy));

        CouponPolicy result = couponPolicyService.findById("123");

        assertThat(result).isNotNull();
        assertThat(result.getId()).isEqualTo("123");
    }

    @Test
    void givenNonExistingId_whenFindById_thenThrowException() {
        when(couponPolicyRepository.findById("999")).thenReturn(Optional.empty());

        assertThrows(CouponPolicyNotFoundException.class, () -> couponPolicyService.findById("999"));
    }

    @Test
    void givenValidCreateRequest_whenCreate_thenReturnSavedCouponPolicy() {
        CouponPolicyDTO.CreateRequest request = CouponPolicyDTO.CreateRequest.builder()
                .name("Summer Sale")
                .description("Discount for summer season")
                .discountType(CouponPolicy.DiscountType.PERCENTAGE)
                .discountValue(10)
                .minimumOrderAmount(100)
                .maximumDiscountAmount(50)
                .totalQuantity(1000)
                .startTime(LocalDateTime.now().minusDays(1))
                .endTime(LocalDateTime.now().plusDays(10))
                .build();

        CouponPolicy savedPolicy = request.toEntity();
        savedPolicy.setId("abc123");

        when(couponPolicyRepository.save(any(CouponPolicy.class))).thenReturn(savedPolicy);

        CouponPolicy result = couponPolicyService.create(request);

        assertThat(result).isNotNull();
        assertThat(result.getId()).isEqualTo("abc123");
        assertThat(result.getName()).isEqualTo("Summer Sale");
        assertThat(result.getDiscountValue()).isEqualTo(10);
    }
}
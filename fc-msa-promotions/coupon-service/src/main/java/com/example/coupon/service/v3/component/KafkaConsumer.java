package com.example.coupon.service.v3.component;

import com.example.coupon.model.CouponDTO;
import com.example.coupon.service.v3.CouponService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
@Slf4j
@RequiredArgsConstructor
public class KafkaConsumer {

    private final CouponService couponService;

    /**
     * Consumes coupon issue requests from Kafka and delegates the actual issuance to the coupon service.
     *
     * @param message the coupon issue request message containing user ID and policy ID
     */
    @KafkaListener(topics = "coupon-issue-requests", groupId = "coupon-service", containerFactory = "couponKafkaListenerContainerFactory")
    public void consumeCouponIssueRequest(CouponDTO.IssueMessage message) {
        try {
            log.info("Received coupon issue request: {}", message);
            couponService.processIssueCoupon(message);
        } catch (Exception e) {
            log.error("Failed to process coupon issue request: {}", e.getMessage(), e);
        }
    }
}

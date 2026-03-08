package com.example.coupon.service.v3.component;

import com.example.coupon.model.CouponDTO;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

@Component
@Slf4j
@RequiredArgsConstructor
public class KafkaProducer {

    private final KafkaTemplate<String, CouponDTO.IssueMessage> kafkaTemplateCouponIssueMessage;

    /**
     * Sends a coupon issue request message to the Kafka topic.
     *
     * @param message the coupon issue message containing policy ID and user ID
     */
    public void sendCouponIssueRequest(CouponDTO.IssueMessage message) {
        String topic = "coupon-issue-requests";
        kafkaTemplateCouponIssueMessage.send(topic, message.getCouponPolicyId(), message)
                .whenComplete((result, ex) -> {
                    if (ex == null) {
                        log.info("Sent message=[{}] with offset=[{}]", message, result.getRecordMetadata().offset());
                    } else {
                        log.error("Unable to send message=[{}] due to : {}", message, ex.getMessage());
                    }
                });
    }
}

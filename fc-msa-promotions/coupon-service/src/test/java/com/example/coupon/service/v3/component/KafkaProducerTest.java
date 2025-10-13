package com.example.coupon.service.v3.component;

import com.example.coupon.model.CouponDTO;
import org.apache.kafka.clients.producer.RecordMetadata;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;

import java.util.concurrent.CompletableFuture;

import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class KafkaProducerTest {

    private static final String TEST_POLICY_ID = "COUPON_POLICY_1";
    private static final String TEST_USER_ID = "USER_1";

    @InjectMocks
    private KafkaProducer kafkaProducer;

    @Mock
    private KafkaTemplate<String, CouponDTO.IssueMessage> kafkaTemplate;

    @Mock
    private SendResult<String, CouponDTO.IssueMessage> sendResult;

    @Mock
    private RecordMetadata recordMetadata;

    private CouponDTO.IssueMessage message;

    @BeforeEach
    void setUp() {
        message = CouponDTO.IssueMessage.builder()
                .couponPolicyId(TEST_POLICY_ID)
                .userId(TEST_USER_ID)
                .build();
    }

    @Test
    @DisplayName("Should send coupon issue request to Kafka topic successfully")
    void sendCouponIssueRequestSuccess() {
        CompletableFuture<SendResult<String, CouponDTO.IssueMessage>> future =
                CompletableFuture.completedFuture(sendResult);

        when(sendResult.getRecordMetadata()).thenReturn(recordMetadata);
        when(recordMetadata.offset()).thenReturn(55L);

        when(kafkaTemplate.send(
                eq("coupon-issue-requests"),
                eq(String.valueOf(message.getCouponPolicyId())),
                eq(message)
        ))
                .thenReturn(future);

        kafkaProducer.sendCouponIssueRequest(message);

        verify(kafkaTemplate, times(1))
                .send("coupon-issue-requests", String.valueOf(message.getCouponPolicyId()), message);

        verify(sendResult, times(1)).getRecordMetadata();
        verify(recordMetadata, times(1)).offset();
    }

    @Test
    @DisplayName("sendCouponIssueRequest → logs error on send failure")
    void sendCouponIssueRequest_logsOnFailure() {
        CompletableFuture<SendResult<String, CouponDTO.IssueMessage>> future = new CompletableFuture<>();
        future.completeExceptionally(new RuntimeException("boom"));

        when(kafkaTemplate.send(anyString(), anyString(), any()))
                .thenReturn(future);

        kafkaProducer.sendCouponIssueRequest(message);

        verify(kafkaTemplate, times(1))
                .send("coupon-issue-requests", String.valueOf(message.getCouponPolicyId()), message);

        // the future is exceptional, so its whenComplete callback runs and hits log.error branch
        // we can't verify the log directly here, but at least we verify the send() call.
    }
}
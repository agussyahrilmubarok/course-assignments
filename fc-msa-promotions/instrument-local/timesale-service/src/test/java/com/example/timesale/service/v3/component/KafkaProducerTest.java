package com.example.timesale.service.v3.component;

import com.example.timesale.model.TimeSaleDTO;
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

    private static final String TEST_REQUEST_ID = "REQ_123";
    private static final String TEST_USER_ID = "USER_1";
    private static final String TEST_TIME_SALE_ID = "SALE_123";

    @InjectMocks
    private KafkaProducer kafkaProducer;

    @Mock
    private KafkaTemplate<String, TimeSaleDTO.PurchaseRequestMessage> kafkaTemplate;

    @Mock
    private SendResult<String, TimeSaleDTO.PurchaseRequestMessage> sendResult;

    @Mock
    private RecordMetadata recordMetadata;

    private TimeSaleDTO.PurchaseRequestMessage message;

    @BeforeEach
    void setUp() {
        message = TimeSaleDTO.PurchaseRequestMessage.builder()
                .userId(TEST_USER_ID)
                .timeSaleId(TEST_TIME_SALE_ID)
                .quantity(2L)
                .build();
    }

    @Test
    @DisplayName("sendPurchaseRequest → should send message successfully")
    void sendPurchaseRequest_success() {
        CompletableFuture<SendResult<String, TimeSaleDTO.PurchaseRequestMessage>> future =
                CompletableFuture.completedFuture(sendResult);

        when(sendResult.getRecordMetadata()).thenReturn(recordMetadata);
        when(recordMetadata.offset()).thenReturn(123L);

        when(kafkaTemplate.send(
                eq("time-sale-requests"),
                eq(TEST_REQUEST_ID),
                eq(message)
        )).thenReturn(future);

        kafkaProducer.sendPurchaseRequest(TEST_REQUEST_ID, message);

        verify(kafkaTemplate, times(1)).send("time-sale-requests", TEST_REQUEST_ID, message);
        verify(sendResult, times(1)).getRecordMetadata();
        verify(recordMetadata, times(1)).offset();
    }

    @Test
    @DisplayName("sendPurchaseRequest → should log error when sending fails")
    void sendPurchaseRequest_failure() {
        CompletableFuture<SendResult<String, TimeSaleDTO.PurchaseRequestMessage>> failedFuture = new CompletableFuture<>();
        failedFuture.completeExceptionally(new RuntimeException("Kafka send failed"));

        when(kafkaTemplate.send(anyString(), anyString(), any())).thenReturn(failedFuture);

        kafkaProducer.sendPurchaseRequest(TEST_REQUEST_ID, message);

        verify(kafkaTemplate).send("time-sale-requests", TEST_REQUEST_ID, message);
        // Logging error can't be directly asserted unless using LogCaptor or a custom logger
    }
}

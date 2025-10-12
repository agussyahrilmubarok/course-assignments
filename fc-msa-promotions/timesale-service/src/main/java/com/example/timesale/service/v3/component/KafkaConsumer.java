package com.example.timesale.service.v3.component;

import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.service.TimeSaleAsyncService;
import com.example.timesale.service.TimeSaleService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
@Slf4j
public class KafkaConsumer {

    private final TimeSaleService timeSaleService;
    private final TimeSaleAsyncService timeSaleAsyncService;

    public KafkaConsumer(@Qualifier("TimeSaleServiceImplV3") TimeSaleService timeSaleService,
                         @Qualifier("TimeSaleServiceImplV3") TimeSaleAsyncService timeSaleAsyncService) {
        this.timeSaleService = timeSaleService;
        this.timeSaleAsyncService = timeSaleAsyncService;
    }

    @KafkaListener(topics = "time-sale-requests", groupId = "timesale-service", containerFactory = "timeSaleKafkaListenerContainerFactory")
    public void consumePurchaseRequest(TimeSaleDTO.PurchaseRequestMessage message) {
        try {
            log.info("Received purchase request: {}", message);
            TimeSaleDTO.PurchaseRequest purchaseRequest = TimeSaleDTO.PurchaseRequest.builder()
                    .timeSaleId(message.getTimeSaleId())
                    .quantity(message.getQuantity())
                    .build();
            timeSaleService.purchase(purchaseRequest, message.getUserId());

            timeSaleAsyncService.savePurchaseResult(message.getRequestId(), "SUCCESS");
        } catch (Exception e) {
            log.error("Failed to process coupon issue request: {}", e.getMessage(), e);
            timeSaleAsyncService.savePurchaseResult(message.getRequestId(), "FAIL");
        } finally {
            timeSaleAsyncService.removePurchaseResultFromQueue(message.getRequestId(), "FAIL");
        }
    }

    @KafkaListener(topics = "time-sale-requests.DLT", groupId = "coupon-dlq-consumer")
    public void consumeFailedPurchaseRequest(TimeSaleDTO.PurchaseRequestMessage message) {
        // Log the failed message for manual investigation or automated alerting
        log.error("Received message from DLQ: {}", message);

        // You can add further logic here, e.g., saving the message to a DB table for retry later,
        // sending notifications, or triggering compensating transactions.
    }

}

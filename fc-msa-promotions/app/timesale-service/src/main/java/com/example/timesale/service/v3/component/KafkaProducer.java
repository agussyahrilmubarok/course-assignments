package com.example.timesale.service.v3.component;

import com.example.timesale.model.TimeSaleDTO;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

@Component
@Slf4j
@RequiredArgsConstructor
public class KafkaProducer {

    private final KafkaTemplate<String, TimeSaleDTO.PurchaseRequestMessage> kafkaTemplateTimeSalePurchaseMessage;

    public void sendPurchaseRequest(String requestId, TimeSaleDTO.PurchaseRequestMessage message) {
        String topic = "time-sale-requests";
        kafkaTemplateTimeSalePurchaseMessage.send(topic, requestId, message)
                .whenComplete((result, ex) -> {
                    if (ex == null) {
                        log.info("Sent message=[{}] with offset=[{}]", message, result.getRecordMetadata().offset());
                    } else {
                        log.error("Unable to send message=[{}] due to : {}", message, ex.getMessage());
                    }
                });
    }
}

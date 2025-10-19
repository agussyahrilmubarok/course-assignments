package com.example.catalog.service;

import com.example.catalog.model.ProductDTO;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

@Component
@Slf4j
@RequiredArgsConstructor
public class EventProducer {

    private final KafkaTemplate<String, byte[]> kafkaTemplate;
    private final ObjectMapper objectMapper;

    public void sendProductTagsAdded(ProductDTO.ProductTagsMessage message) {
        String topic = "product-tags-added";
        try {
            byte[] serializedMessage = objectMapper.writeValueAsBytes(message);

            kafkaTemplate.send(topic, message.getProductId(), serializedMessage)
                    .whenComplete((result, ex) -> {
                        if (ex == null) {
                            log.info("Sent message [{}] to topic [{}] with offset [{}]", message, topic, result.getRecordMetadata().offset());
                        } else {
                            log.error("Failed to send message [{}] to topic [{}]: {}", message, topic, ex.getMessage(), ex);
                        }
                    });
        } catch (Exception e) {
            log.error("Failed to serialize message [{}]: {}", message, e.getMessage(), e);
        }
    }

    public void sendProductTagsRemoved(ProductDTO.ProductTagsMessage message) {
        String topic = "product-tags-removed";
        try {
            byte[] serializedMessage = objectMapper.writeValueAsBytes(message);

            kafkaTemplate.send(topic, message.getProductId(), serializedMessage)
                    .whenComplete((result, ex) -> {
                        if (ex == null) {
                            log.info("Sent message [{}] to topic [{}] with offset [{}]", message, topic, result.getRecordMetadata().offset());
                        } else {
                            log.error("Failed to send message [{}] to topic [{}]: {}", message, topic, ex.getMessage(), ex);
                        }
                    });
        } catch (Exception e) {
            log.error("Failed to serialize message [{}]: {}", message, e.getMessage(), e);
        }
    }
}


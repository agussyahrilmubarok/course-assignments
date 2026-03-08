package com.example.search.service;

import com.example.search.dto.ProductDTO;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
@Slf4j
@RequiredArgsConstructor
public class EventListener {

    private final ObjectMapper objectMapper;
    private final SearchService searchService;

    @KafkaListener(
            topics = "product-tags-added",
            groupId = "${spring.kafka.consumer.group-id}",
            containerFactory = "kafkaListenerContainerFactory"
    )
    public void consumeProductTagsAdded(ConsumerRecord<String, byte[]> record) {
        try {
            byte[] value = record.value();
            ProductDTO.ProductTagsMessage message = objectMapper.readValue(value, ProductDTO.ProductTagsMessage.class);
            log.info("Received message from topic [{}]: {}", record.topic(), message);
            searchService.addTagsCache(message.getProductId(), message.getTags());
        } catch (Exception e) {
            log.error("Failed to deserialize Kafka message: {}", e.getMessage(), e);
        }
    }

    @KafkaListener(
            topics = "product-tags-removed",
            groupId = "${spring.kafka.consumer.group-id}",
            containerFactory = "kafkaListenerContainerFactory"
    )
    public void consumeProductTagsRemoved(ConsumerRecord<String, byte[]> record) {
        try {
            byte[] value = record.value();
            ProductDTO.ProductTagsMessage message = objectMapper.readValue(value, ProductDTO.ProductTagsMessage.class);
            log.info("Received message from topic [{}]: {}", record.topic(), message);
            searchService.removeTagsCache(message.getProductId(), message.getTags());
        } catch (Exception e) {
            log.error("Failed to deserialize Kafka message: {}", e.getMessage(), e);
        }
    }
}

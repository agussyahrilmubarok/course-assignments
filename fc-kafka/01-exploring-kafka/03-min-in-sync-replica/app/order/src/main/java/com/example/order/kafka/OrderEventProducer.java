package com.example.order.kafka;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.common.header.Header;
import org.apache.kafka.common.header.internals.RecordHeader;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Component;

import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;

@Component
@Slf4j
public class OrderEventProducer {

    @Value("${spring.kafka.topics.order-events.name}")
    private String orderEventsTopicName;

    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;

    public OrderEventProducer(KafkaTemplate<String, String> kafkaTemplate, ObjectMapper objectMapper) {
        this.kafkaTemplate = kafkaTemplate;
        this.objectMapper = objectMapper;
    }

    /**
     * Serialize OrderEvent and send it asynchronously to Kafka.
     */
    public CompletableFuture<SendResult<String, String>> sendEvent(OrderEvent event) throws JsonProcessingException {
        String key = event.getOrder().getOrderId(); // partition key based on orderId
        String value = objectMapper.writeValueAsString(event);

        ProducerRecord<String, String> record = buildRecord(key, value, orderEventsTopicName);

        CompletableFuture<SendResult<String, String>> future = kafkaTemplate.send(record);

        future.whenComplete((result, ex) -> {
            if (ex != null) {
                handleFailure(key, value, ex);
            } else {
                handleSuccess(key, value, result);
            }
        });

        return future;
    }

    /**
     * Build Kafka ProducerRecord with topic, key, and value.
     */
    private ProducerRecord<String, String> buildRecord(String key, String value, String topic) {
        List<Header> headers = List.of(
                new RecordHeader("event-type", "OrderEvent".getBytes(StandardCharsets.UTF_8)),
                new RecordHeader("event-source", "order-service".getBytes(StandardCharsets.UTF_8)),
                new RecordHeader("correlation-id", UUID.randomUUID().toString().getBytes(StandardCharsets.UTF_8)),
                new RecordHeader("trace-id", UUID.randomUUID().toString().getBytes(StandardCharsets.UTF_8)),
                new RecordHeader("timestamp", String.valueOf(System.currentTimeMillis()).getBytes(StandardCharsets.UTF_8)),
                new RecordHeader("content-type", "application/json".getBytes(StandardCharsets.UTF_8)),
                new RecordHeader("initiator", "system-user".getBytes(StandardCharsets.UTF_8))
        );

        return new ProducerRecord<>(topic, null, key, value, headers);
    }

    /**
     * Handle failure case when sending message to Kafka fails.
     */
    private void handleFailure(String key, String value, Throwable throwable) {
        log.error("Failed to send OrderEvent. key={}, value={}, error={}", key, value, throwable.getMessage(), throwable);
    }

    /**
     * Handle success case when Kafka acknowledges the message.
     */
    private void handleSuccess(String key, String value, SendResult<String, String> sendResult) {
        log.info("Successfully sent OrderEvent. key={}, value={}, topic={}, partition={}, offset={}",
                key,
                value,
                sendResult.getRecordMetadata().topic(),
                sendResult.getRecordMetadata().partition(),
                sendResult.getRecordMetadata().offset()
        );
    }
}

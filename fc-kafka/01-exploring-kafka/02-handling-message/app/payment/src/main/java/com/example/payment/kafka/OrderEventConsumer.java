package com.example.payment.kafka;

import com.example.payment.domain.Payment;
import com.example.payment.model.PaymentDTO;
import com.example.payment.service.PaymentService;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.consumer.Consumer;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.OffsetAndMetadata;
import org.apache.kafka.common.TopicPartition;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import java.util.Collections;

@Component
@Slf4j
public class OrderEventConsumer {

    private final ObjectMapper objectMapper;
    private final PaymentService paymentService;

    public OrderEventConsumer(ObjectMapper objectMapper, PaymentService paymentService) {
        this.objectMapper = objectMapper;
        this.paymentService = paymentService;
    }

    @KafkaListener(
            topics = "${spring.kafka.topic.order.events}",
            groupId = "${spring.kafka.consumer.group-id:payment-service-consumer}",
            containerFactory = "kafkaListenerContainerFactory"
    )
    public void consume(ConsumerRecord<String, String> record, Consumer<String, String> consumer) {
        try {
            processOrderEvent(record);
            TopicPartition topicPartition = new TopicPartition(record.topic(), record.partition());
            OffsetAndMetadata offsetAndMetadata = new OffsetAndMetadata(record.offset() + 1);
            consumer.commitSync(Collections.singletonMap(topicPartition, offsetAndMetadata));
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }

    private void processOrderEvent(ConsumerRecord<String, String> consumerRecord) throws JsonProcessingException {
        OrderEvent event = objectMapper.readValue(consumerRecord.value(), OrderEvent.class);

        switch (event.getEventType()) {
            case OrderEvent.OrderEventType.CREATED:
                PaymentDTO paymentDTO = new PaymentDTO();
                paymentDTO.setOrderId(event.getOrder().getOrderId());
                paymentDTO.setCustomerId(event.getOrder().getCustomerId());
                paymentDTO.setAmount(event.getOrder().getTotalAmount());
                paymentDTO.setStatus(Payment.Status.PENDING.name());
                paymentService.create(paymentDTO);
                break;
            default:
                log.info("Invalid Order Event Type");
        }
    }
}

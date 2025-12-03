package com.example.order.kafka;

import org.apache.kafka.clients.admin.NewTopic;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.TopicBuilder;

@Configuration
public class KafkaConfig {

    @Value("${spring.kafka.topics.order-events.name}")
    private String orderEventsTopicName;

    @Value("${spring.kafka.topics.order-events.partitions}")
    private int orderEventsPartitions;

    @Value("${spring.kafka.topics.order-events.replicas}")
    private short orderEventsReplicas;

    @Value("${spring.kafka.topics.order-events.min-insync-replicas}")
    private String orderEventsMinInSyncReplicas;

    @Bean
    public NewTopic orderEventsTopic() {
        return TopicBuilder.name(orderEventsTopicName)
                .partitions(orderEventsPartitions)
                .replicas(orderEventsReplicas)
                .config("min.insync.replicas", orderEventsMinInSyncReplicas)
                .build();
    }
}

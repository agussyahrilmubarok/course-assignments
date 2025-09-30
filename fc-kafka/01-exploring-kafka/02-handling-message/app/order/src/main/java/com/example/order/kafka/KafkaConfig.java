package com.example.order.kafka;

import org.apache.kafka.clients.admin.NewTopic;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.TopicBuilder;

@Configuration
public class KafkaConfig {

    @Value("${spring.kafka.topic.order.events}")
    private String orderEventsTopic;

    @Bean
    public NewTopic orderTopic() {
        //
        return TopicBuilder.name(orderEventsTopic)
                .partitions(3)
                .replicas(3)
                .build();
    }
}

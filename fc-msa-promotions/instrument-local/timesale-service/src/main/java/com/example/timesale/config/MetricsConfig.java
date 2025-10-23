package com.example.timesale.config;

import com.example.timesale.aop.TimeSaleMetricsAspect;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Metrics;
import io.micrometer.core.instrument.binder.MeterBinder;
import lombok.RequiredArgsConstructor;
import org.redisson.api.RedissonClient;
import org.redisson.api.redisnode.RedisNode;
import org.redisson.api.redisnode.RedisNodes;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.Map;

@Configuration
@RequiredArgsConstructor
public class MetricsConfig {

    @Bean
    public TimeSaleMetricsAspect timeSaleMetricsAspect(MeterRegistry registry) {
        return new TimeSaleMetricsAspect(registry);
    }

    @Bean
    public MeterBinder redisMetrics(RedissonClient redissonClient) {
        return registry -> {
            Metrics.gauge("redis.memory_used_bytes", redissonClient, client -> {
                try {
                    Map<String, String> info = client.getRedisNodes(RedisNodes.SINGLE).getInstance().info(RedisNode.InfoSection.MEMORY);
                    return Double.parseDouble(info.get("used_memory"));
                } catch (Exception e) {
                    return 0.0;
                }
            });

            Metrics.gauge("redis.connected_clients", redissonClient, client -> {
                try {
                    Map<String, String> info = client.getRedisNodes(RedisNodes.SINGLE).getInstance().info(RedisNode.InfoSection.CLIENTS);
                    return Double.parseDouble(info.get("connected_clients"));
                } catch (Exception e) {
                    return 0.0;
                }
            });
        };
    }
}

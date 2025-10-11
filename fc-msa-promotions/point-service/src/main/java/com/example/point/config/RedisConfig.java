package com.example.point.config;

import org.redisson.Redisson;
import org.redisson.api.RedissonClient;
import org.redisson.client.codec.StringCodec;
import org.redisson.config.Config;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RedisConfig {

    @Value("${spring.data.redis.host}")
    private String host;

    @Value("${spring.data.redis.port}")
    private int port;

    @Bean
    public RedissonClient redissonClient() {
        Config config = new Config();
        config.setCodec(new StringCodec()); // Use it for java.lang.ClassCastException: class java.lang.String cannot be cast to class java.lang.Long
        // config.setCodec(new JsonJacksonCodec());
        config.useSingleServer().setAddress("redis://" + host + ":" + port);

        return Redisson.create(config);
    }
}
package com.example.point.scheduler;

import com.example.point.domain.Point;
import com.example.point.domain.PointBalance;
import com.example.point.repos.PointBalanceRepository;
import com.example.point.repos.PointRepository;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import org.redisson.api.RBucket;
import org.redisson.api.RedissonClient;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Duration;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.List;

@Service
@RequiredArgsConstructor
public class RedisSyncService {

    private static final String POINT_CACHE_PREFIX = "cache:point:";
    private static final String BALANCE_CACHE_PREFIX = "cache:balance:";

    private final PointRepository pointRepository;
    private final PointBalanceRepository pointBalanceRepository;
    private final RedissonClient redissonClient;
    private final ObjectMapper objectMapper;

    @Scheduled(fixedRate = 5 * 60 * 1000)
    public void syncDataToRedis() throws JsonProcessingException {
        Duration ttl = calculateTTLUntil2AM();
        syncPoints(ttl);
        syncPointBalances(ttl);
    }

    private void syncPoints(Duration ttl) throws JsonProcessingException {
        List<Point> points = pointRepository.findAll();
        for (Point point : points) {
            String key = POINT_CACHE_PREFIX + point.getId();
            PointRedisDTO dto = PointRedisDTO.toDto(point);
            String pointJson = objectMapper.writeValueAsString(dto);
            RBucket<String> bucket = redissonClient.getBucket(key);
            bucket.set(pointJson, ttl);
        }
    }

    @Transactional
    private void syncPointBalances(Duration ttl) throws JsonProcessingException {
        ObjectMapper mapper = new ObjectMapper();
        List<PointBalance> balances = pointBalanceRepository.findAll();
        for (PointBalance balance : balances) {
            String key = BALANCE_CACHE_PREFIX + balance.getId();
            PointBalanceRedisDTO dto = PointBalanceRedisDTO.toDto(balance);
            String balanceJson = objectMapper.writeValueAsString(dto);
            RBucket<String> bucket = redissonClient.getBucket(key);
            bucket.set(balanceJson, ttl);
        }
    }

    @Transactional
    private Duration calculateTTLUntil2AM() {
        LocalDateTime now = LocalDateTime.now();
        LocalDateTime next2AM = now.toLocalTime().isBefore(LocalTime.of(2, 0))
                ? now.with(LocalTime.of(2, 0))
                : now.plusDays(1).with(LocalTime.of(2, 0));

        return Duration.between(now, next2AM);
    }
}

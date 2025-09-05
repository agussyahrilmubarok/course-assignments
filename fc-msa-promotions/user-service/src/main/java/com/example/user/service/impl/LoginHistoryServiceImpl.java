package com.example.user.service.impl;

import com.example.user.model.LoginHistory;
import com.example.user.service.LoginHistoryService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Service
@Slf4j
@RequiredArgsConstructor
public class LoginHistoryServiceImpl implements LoginHistoryService {

    private final RedisTemplate<String, Object> redisTemplate;

    @Override
    public void recordLogin(String userId, String ipAddress) {
        LoginHistory history = new LoginHistory();
        history.setUserId(userId);
        history.setIpAddress(ipAddress);
        history.setLoginAt(LocalDateTime.now());

        String key = "users:login:history:" + userId;
        redisTemplate.opsForList().leftPush(key, history);
        redisTemplate.expire(key, java.time.Duration.ofDays(30));

        log.info("Recorded login for userId={} from ipAddress={} at {}",
                userId, ipAddress, history.getLoginAt());
    }
}

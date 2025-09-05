package com.example.user.service.impl;

import com.example.user.model.LoginHistory;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.*;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.redis.core.ListOperations;
import org.springframework.data.redis.core.RedisTemplate;

import java.time.Duration;

import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class LoginHistoryServiceImplTest {

    @Mock
    private RedisTemplate<String, Object> redisTemplate;

    @Mock
    private ListOperations<String, Object> listOperations;

    @InjectMocks
    private LoginHistoryServiceImpl loginHistoryService;

    @Test
    void testRecordLogin() {
        String userId = "user123";
        String ipAddress = "192.168.1.10";
        String expectedKey = "users:login:history:" + userId;

        when(redisTemplate.opsForList()).thenReturn(listOperations);

        loginHistoryService.recordLogin(userId, ipAddress);

        verify(listOperations, times(1)).leftPush(eq(expectedKey), any(LoginHistory.class));
        verify(redisTemplate, times(1)).expire(eq(expectedKey), eq(Duration.ofDays(30)));
    }
}

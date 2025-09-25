package com.example.user.service.impl;

import com.auth0.jwt.interfaces.DecodedJWT;
import com.example.user.domain.User;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.test.util.ReflectionTestUtils;

import static org.junit.jupiter.api.Assertions.*;

@ExtendWith(MockitoExtension.class)
class JWTServiceImplTest {

    @InjectMocks
    private JWTServiceImpl jwtService;

    @BeforeEach
    void setUp() {
        ReflectionTestUtils.setField(jwtService, "jwtSecretKey", "testSecretKeyForJWT1234567890");
        ReflectionTestUtils.setField(jwtService, "jwtExpirationMs", 3600000);
    }

    @Test
    void testGenerateAndValidateToken() {
        User user = new User();
        user.setEmail("test@example.com");

        String token = jwtService.generate(user);

        assertNotNull(token, "Token should not be null");
        assertFalse(token.isEmpty(), "Token should not be empty");

        DecodedJWT decoded = jwtService.validateToken(token);

        assertEquals("test@example.com", decoded.getSubject());
        assertEquals("USER", decoded.getClaim("role").asString());
        assertNotNull(decoded.getIssuedAt());
        assertNotNull(decoded.getExpiresAt());
    }

    @Test
    void testValidateToken_withInvalidToken_shouldThrowException() {
        String invalidToken = "invalid.jwt.token";

        Exception exception = assertThrows(IllegalArgumentException.class,
                () -> jwtService.validateToken(invalidToken));

        assertEquals("Invalid token", exception.getMessage());
    }
}

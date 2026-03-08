package com.example.member.service;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.auth0.jwt.interfaces.JWTVerifier;
import com.example.member.domain.User;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.Date;

@Service
@Slf4j
@RequiredArgsConstructor
public class JWTServiceImpl implements JWTService {

    private final PasswordEncoder passwordEncoder;

    @Value("${jwt.secret.key}")
    private String jwtSecretKey;
    @Value("${jwt.expires}")
    private long jwtExpirationMs;

    @Override
    public String generate(User user) {
        long now = System.currentTimeMillis();

        return JWT.create()
                .withSubject(user.getUsername())
                .withClaim("role", "USER")
                .withIssuedAt(new Date(now))
                .withExpiresAt(new Date(now + jwtExpirationMs))
                .sign(getAlgorithm());
    }

    @Override
    public DecodedJWT validateToken(String token) {
        try {
            DecodedJWT decodedJWT = getVerifier().verify(token);
            log.info("JWT token is valid for user: {}", decodedJWT.getSubject());
            return decodedJWT;
        } catch (Exception e) {
            log.error("Invalid JWT token", e);
            throw new IllegalArgumentException("Invalid token");
        }
    }

    private Algorithm getAlgorithm() {
        return Algorithm.HMAC512(jwtSecretKey);
    }

    private JWTVerifier getVerifier() {
        return JWT.require(getAlgorithm()).build();
    }
}

package com.example.witrack.backend.security.jwt;

import com.auth0.jwt.JWT;
import com.auth0.jwt.JWTVerifier;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.security.user.UserDetailsImpl;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Component;

import java.util.Date;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

@Component
public class JwtProvider {

    @Value("${jwt.secret_key}")
    private String jwtSecretKey;

    @Value("${jwt.expiry_in}")
    private long jwtExpiryIn;

    public String generateToken(Authentication authentication) {
        UserDetailsImpl userDetails = (UserDetailsImpl) authentication.getPrincipal();

        String userId = userDetails.getId().toString();
        List<String> roles = userDetails.getAuthorities()
                .stream()
                .map(auth -> auth.getAuthority())
                .collect(Collectors.toList());

        return createToken(userId, roles);
    }

    public String generateToken(String subject, Set<User.Role> roles) {
        List<String> roleNames = roles.stream()
                .map(User.Role::name)
                .collect(Collectors.toList());
        return createToken(subject, roleNames);
    }

    public boolean validateToken(String token) {
        try {
            decodeToken(token);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    public String extractSubject(String token) {
        return decodeToken(token).getSubject();
    }

    public List<String> extractRoles(String token) {
        return decodeToken(token)
                .getClaim("roles")
                .asList(String.class);
    }

    private String createToken(String subject, List<String> roles) {
        Date now = new Date();
        Date expiryDate = new Date(now.getTime() + jwtExpiryIn);

        return JWT.create()
                .withSubject(subject)
                .withIssuedAt(now)
                .withExpiresAt(expiryDate)
                .withClaim("roles", roles)
                .sign(getAlgorithm());
    }

    private DecodedJWT decodeToken(String token) {
        if (token == null || token.isBlank()) {
            throw new RuntimeException("JWT token is null or empty");
        }

        if (token.startsWith("Bearer ")) {
            token = token.substring(7);
        }

        return getVerifier().verify(token);
    }

    private Algorithm getAlgorithm() {
        return Algorithm.HMAC512(jwtSecretKey);
    }

    private JWTVerifier getVerifier() {
        return JWT.require(getAlgorithm()).build();
    }
}

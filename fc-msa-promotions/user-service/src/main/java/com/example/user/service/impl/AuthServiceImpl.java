package com.example.user.service.impl;

import com.auth0.jwt.interfaces.DecodedJWT;
import com.example.user.domain.User;
import com.example.user.exception.DuplicateUserException;
import com.example.user.exception.UnauthorizedAccessException;
import com.example.user.exception.UserNotFoundException;
import com.example.user.model.AuthDTO;
import com.example.user.model.UserDTO;
import com.example.user.repos.UserRepository;
import com.example.user.service.AuthService;
import com.example.user.service.JWTService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class AuthServiceImpl implements AuthService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final JWTService jwtService;

    @Override
    public void signUp(AuthDTO.SignUp param) {
        if (userRepository.existsByEmailIgnoreCase(param.getEmail())) {
            log.error("Duplicate user registration attempt for email={}", param.getEmail());
            throw new DuplicateUserException("User already exists with email " + param.getEmail());
        }

        User user = new User();
        user.setId(UUID.randomUUID().toString());
        user.setName(param.getName().toUpperCase());
        user.setEmail(param.getEmail().toUpperCase());
        user.setPassword(passwordEncoder.encode(param.getPassword()));
        userRepository.save(user);
        log.info("User successfully registered with id={} email={}", user.getId(), user.getEmail());
    }

    @Override
    public AuthDTO.Response signIn(AuthDTO.SignIn param) {
        User user = userRepository.findByEmail(param.getEmail().toUpperCase())
                .orElseThrow(() -> {
                    log.error("Login failed: user not found with email={}", param.getEmail());
                    return new UserNotFoundException("User not found with email " + param.getEmail());
                });

        if (!passwordEncoder.matches(param.getPassword(), user.getPassword())) {
            log.error("Login failed: invalid password for email={}", param.getEmail());
            throw new UnauthorizedAccessException("Invalid password");
        }

        return AuthDTO.Response.builder()
                .token(jwtService.generate(user))
                .user(UserDTO.from(user))
                .build();
    }

    @Override
    public AuthDTO.Response validateToken(AuthDTO.TokenRequest param) {
        DecodedJWT claims = jwtService.validateToken(param.getToken());

        User user = userRepository.findByEmail(claims.getSubject().toUpperCase())
                .orElseThrow(() -> {
                    log.error("Token validation failed: user not found with email={}", claims.getSubject());
                    return new UserNotFoundException("User not found with email " + claims.getSubject());
                });

        return AuthDTO.Response.builder()
                .token(param.getToken())
                .user(UserDTO.from(user))
                .build();
    }
}

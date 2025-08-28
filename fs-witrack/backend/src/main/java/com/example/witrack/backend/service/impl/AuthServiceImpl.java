package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.DuplicateFieldException;
import com.example.witrack.backend.model.*;
import com.example.witrack.backend.repository.UserRepository;
import com.example.witrack.backend.security.jwt.JwtProvider;
import com.example.witrack.backend.service.AuthService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class AuthServiceImpl implements AuthService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final AuthenticationManager authenticationManager;
    private final JwtProvider jwtProvider;

    @Override
    public SignUpResponse signUp(SignUpRequest request) {
        String email = request.getEmail().toLowerCase();
        if (userRepository.existsByEmailIgnoreCase(email)) {
            log.warn("Sign Up failed: email {} is already in use", email);
            throw new DuplicateFieldException("email", "Email is already in use");
        }

        User user = new User();
        user.setId(UUID.randomUUID().toString());
        user.setFullName(request.getFullName());
        user.setEmail(email);
        user.setPassword(passwordEncoder.encode(request.getPassword()));
        user.setRoles(Collections.singleton(User.Role.ROLE_USER));

        User savedUser = userRepository.save(user);
        String token = authenticateAndGenerateToken(email, request.getPassword());

        log.info("User {} successfully registered with id {}", email, savedUser.getId());
        return SignUpResponse.builder()
                .token(token)
                .user(UserResponse.fromUser(savedUser))
                .build();
    }

    @Override
    public SignInResponse signIn(SignInRequest request) {
        String email = request.getEmail().toLowerCase();
        String token = authenticateAndGenerateToken(email, request.getPassword());

        User user = userRepository.findByEmail(email)
                .orElseThrow(() -> new UsernameNotFoundException("User is not found"));

        log.info("User {} successfully authenticated", email);
        return SignInResponse.builder()
                .token(token)
                .user(UserResponse.fromUser(user))
                .build();
    }

    private String authenticateAndGenerateToken(String email, String password) {
        UsernamePasswordAuthenticationToken authToken =
                new UsernamePasswordAuthenticationToken(email, password);

        try {
            Authentication authentication = authenticationManager.authenticate(authToken);
            SecurityContextHolder.getContext().setAuthentication(authentication);
            log.debug("Authentication successful for user {}", email);
            return jwtProvider.generateToken(authentication);
        } catch (BadCredentialsException e) {
            log.warn("Authentication failed: bad credentials for email {}", email);
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during authentication for email {}: {}", email, e.getMessage());
            throw e;
        }
    }
}

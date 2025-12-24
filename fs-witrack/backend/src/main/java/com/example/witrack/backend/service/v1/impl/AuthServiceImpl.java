package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.DuplicateFieldException;
import com.example.witrack.backend.model.AuthDTO;
import com.example.witrack.backend.model.UserDTO;
import com.example.witrack.backend.repos.UserRepository;
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

@Service
@Slf4j
@RequiredArgsConstructor
public class AuthServiceImpl implements AuthService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtProvider jwtProvider;
    private final AuthenticationManager authenticationManager;

    @Override
    public AuthDTO.AuthResponse signUp(AuthDTO.SignUpRequest request) {
        String email = request.getEmail().toLowerCase();
        if (userRepository.existsByEmailIgnoreCase(email)) {
            log.warn("Sign up failed: email already exists, email={}", email);
            throw new DuplicateFieldException("email", "Email is already in use");
        }

        User user = new User();
        user.setFullName(request.getFullName());
        user.setEmail(email);
        user.setPassword(passwordEncoder.encode(request.getPassword()));
        user.setRoles(Collections.singleton(User.Role.ROLE_USER));

        user = userRepository.save(user);
        String token = jwtProvider.generateToken(user.getId().toString(), user.getRoles());

        log.info("User registered successfully, userId={}, email={}", user.getId(), email);
        return AuthDTO.AuthResponse.builder()
                .token(token)
                .user(UserDTO.UserResponse.fromUser(user))
                .build();
    }

    @Override
    public AuthDTO.AuthResponse signIn(AuthDTO.SignInRequest request) {
        String email = request.getEmail().toLowerCase();
        String token = authenticateAndGenerateToken(email, request.getPassword());

        User user = userRepository.findByEmail(email)
                .orElseThrow(() -> {
                    log.warn("Sign in failed: user not found, email={}", email);
                    return new UsernameNotFoundException("User is not found");
                });

        log.info("Sign in successful, userId={}, email={}", user.getId(), email);
        return AuthDTO.AuthResponse.builder()
                .token(token)
                .user(UserDTO.UserResponse.fromUser(user))
                .build();
    }

    private String authenticateAndGenerateToken(String email, String password) {
        UsernamePasswordAuthenticationToken authToken = new UsernamePasswordAuthenticationToken(email, password);

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

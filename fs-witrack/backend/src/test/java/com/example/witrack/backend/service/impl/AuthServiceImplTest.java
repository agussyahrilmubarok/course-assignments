package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.DuplicateFieldException;
import com.example.witrack.backend.model.SignInRequest;
import com.example.witrack.backend.model.SignInResponse;
import com.example.witrack.backend.model.SignUpRequest;
import com.example.witrack.backend.model.SignUpResponse;
import com.example.witrack.backend.repository.UserRepository;
import com.example.witrack.backend.security.jwt.JwtProvider;
import com.example.witrack.backend.service.BaseServiceTest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.util.Collections;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

class AuthServiceImplTest extends BaseServiceTest {

    @InjectMocks
    private AuthServiceImpl authService;

    @Mock
    private UserRepository userRepository;

    @Mock
    private PasswordEncoder passwordEncoder;

    @Mock
    private AuthenticationManager authenticationManager;

    @Mock
    private JwtProvider jwtProvider;

    @Mock
    private Authentication authentication;

    private SignUpRequest signUpRequest;
    private SignInRequest signInRequest;
    private User testUser;

    @BeforeEach
    void setUp() {
        signUpRequest = SignUpRequest.builder()
                .fullName("John Doe")
                .email("john@example.com")
                .password("password123")
                .build();

        signInRequest = SignInRequest.builder()
                .email("john@example.com")
                .password("password123")
                .build();

        testUser = new User();
        testUser.setId(UUID.randomUUID().toString());
        testUser.setFullName(signUpRequest.getFullName());
        testUser.setEmail(signUpRequest.getEmail().toLowerCase());
        testUser.setPassword("encodedPassword");
        testUser.setRoles(Collections.singleton(User.Role.ROLE_USER));
    }

    @Test
    void givenValidRequest_whenSignUp_thenReturnSignUpResponse() {
        when(userRepository.existsByEmailIgnoreCase(anyString())).thenReturn(false);
        when(passwordEncoder.encode(anyString())).thenReturn("encodedPassword");
        when(userRepository.save(any(User.class))).thenReturn(testUser);
        when(authenticationManager.authenticate(any(UsernamePasswordAuthenticationToken.class)))
                .thenReturn(authentication);
        when(jwtProvider.generateToken(authentication)).thenReturn("jwt-token");

        SignUpResponse response = authService.signUp(signUpRequest);

        assertNotNull(response);
        assertEquals("jwt-token", response.getToken());
        assertEquals(testUser.getEmail(), response.getUser().getEmail());

        verify(userRepository, times(1)).save(any(User.class));
    }

    @Test
    void givenExistingEmail_whenSignUp_thenThrowDuplicateFieldException() {
        when(userRepository.existsByEmailIgnoreCase(anyString())).thenReturn(true);

        DuplicateFieldException exception = assertThrows(DuplicateFieldException.class,
                () -> authService.signUp(signUpRequest));

        assertEquals("Email is already in use", exception.getMessage());
        assertEquals("email", exception.getField());
        verify(userRepository, never()).save(any(User.class));
    }

    @Test
    void givenValidCredentials_whenSignIn_thenReturnSignInResponse() {
        when(authenticationManager.authenticate(any(UsernamePasswordAuthenticationToken.class)))
                .thenReturn(authentication);
        when(jwtProvider.generateToken(authentication)).thenReturn("jwt-token");
        when(userRepository.findByEmail(anyString())).thenReturn(Optional.of(testUser));

        SignInResponse response = authService.signIn(signInRequest);

        assertNotNull(response);
        assertEquals("jwt-token", response.getToken());
        assertEquals(testUser.getEmail(), response.getUser().getEmail());
    }

    @Test
    void givenNonExistingEmail_whenSignIn_thenThrowUsernameNotFoundException() {
        when(authenticationManager.authenticate(any(UsernamePasswordAuthenticationToken.class)))
                .thenReturn(authentication);
        when(jwtProvider.generateToken(authentication)).thenReturn("jwt-token");
        when(userRepository.findByEmail(anyString())).thenReturn(Optional.empty());

        UsernameNotFoundException exception = assertThrows(UsernameNotFoundException.class,
                () -> authService.signIn(signInRequest));

        assertTrue(exception.getMessage().contains("User is not found"));
    }

    @Test
    void givenInvalidCredentials_whenSignIn_thenThrowBadCredentialsException() {
        when(authenticationManager.authenticate(any(UsernamePasswordAuthenticationToken.class)))
                .thenThrow(new BadCredentialsException("Bad credentials"));

        assertThrows(BadCredentialsException.class, () -> authService.signIn(signInRequest));
    }

    @Test
    void givenUnexpectedError_whenAuthenticate_thenThrowException() {
        when(authenticationManager.authenticate(any(UsernamePasswordAuthenticationToken.class)))
                .thenThrow(new RuntimeException("Unexpected error"));

        RuntimeException exception = assertThrows(RuntimeException.class,
                () -> authService.signIn(signInRequest));

        assertEquals("Unexpected error", exception.getMessage());
    }
}
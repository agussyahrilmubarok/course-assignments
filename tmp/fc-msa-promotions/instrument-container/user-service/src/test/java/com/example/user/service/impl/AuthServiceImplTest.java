package com.example.user.service.impl;

import com.auth0.jwt.interfaces.DecodedJWT;
import com.example.user.domain.User;
import com.example.user.exception.DuplicateUserException;
import com.example.user.exception.UnauthorizedAccessException;
import com.example.user.exception.UserNotFoundException;
import com.example.user.model.AuthDTO;
import com.example.user.model.UserDTO;
import com.example.user.repos.UserRepository;
import com.example.user.service.JWTService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AuthServiceImplTest {

    @InjectMocks
    private AuthServiceImpl authService;

    @Mock
    private UserRepository userRepository;
    @Mock
    private PasswordEncoder passwordEncoder;
    @Mock
    private JWTService jwtService;

    private User user;

    @BeforeEach
    void setUp() {
        user = new User();
        user.setId("123");
        user.setName("Test User");
        user.setEmail("TEST@EXAMPLE.COM");
        user.setPassword("encodedPassword");
    }

    @Test
    void testSignUp_whenValidRequest_shouldSaveUser() {
        AuthDTO.SignUp signUp = AuthDTO.SignUp.builder()
                .name("Test")
                .email("test@example.com")
                .password("password")
                .build();

        when(userRepository.existsByEmailIgnoreCase(signUp.getEmail())).thenReturn(false);
        when(passwordEncoder.encode(signUp.getPassword())).thenReturn("encodedPassword");

        authService.signUp(signUp);

        verify(userRepository, times(1)).save(any(User.class));
    }

    @Test
    void testSignUp_whenEmailExists_shouldThrowDuplicateUserException() {
        AuthDTO.SignUp signUp = AuthDTO.SignUp.builder()
                .name("Test")
                .email("test@example.com")
                .password("password")
                .build();

        when(userRepository.existsByEmailIgnoreCase(signUp.getEmail())).thenReturn(true);

        assertThrows(DuplicateUserException.class, () -> authService.signUp(signUp));

        verify(userRepository, never()).save(any(User.class));
    }

    @Test
    void testSignIn_whenValidCredentials_shouldReturnTokenAndUser() {
        AuthDTO.SignIn signIn = AuthDTO.SignIn.builder()
                .email("test@example.com")
                .password("password")
                .build();

        when(userRepository.findByEmail("TEST@EXAMPLE.COM")).thenReturn(Optional.of(user));
        when(passwordEncoder.matches("password", "encodedPassword")).thenReturn(true);
        when(jwtService.generate(user)).thenReturn("mockedToken");

        AuthDTO.Response response = authService.signIn(signIn);

        assertNotNull(response);
        assertEquals("mockedToken", response.getToken());
        assertEquals(UserDTO.from(user), response.getUser());
    }

    @Test
    void testSignIn_whenUserNotFound_shouldThrowUserNotFoundException() {
        AuthDTO.SignIn signIn = AuthDTO.SignIn.builder()
                .email("notfound@example.com")
                .password("password")
                .build();

        when(userRepository.findByEmail("NOTFOUND@EXAMPLE.COM")).thenReturn(Optional.empty());

        assertThrows(UserNotFoundException.class, () -> authService.signIn(signIn));

        verify(jwtService, never()).generate(any());
    }

    @Test
    void testSignIn_whenInvalidPassword_shouldThrowUnauthorizedAccessException() {
        AuthDTO.SignIn signIn = AuthDTO.SignIn.builder()
                .email("test@example.com")
                .password("wrongPassword")
                .build();

        when(userRepository.findByEmail("TEST@EXAMPLE.COM")).thenReturn(Optional.of(user));
        when(passwordEncoder.matches("wrongPassword", "encodedPassword")).thenReturn(false); // ❌ password tidak cocok

        Exception ex = assertThrows(UnauthorizedAccessException.class, () -> authService.signIn(signIn));
        assertEquals("Invalid password", ex.getMessage());

        verify(jwtService, never()).generate(any());
    }

    @Test
    void testValidateToken_whenValidToken_shouldReturnUser() {
        AuthDTO.TokenRequest request = AuthDTO.TokenRequest.builder()
                .token("validToken")
                .build();

        DecodedJWT decodedJWT = mock(DecodedJWT.class);
        when(decodedJWT.getSubject()).thenReturn("test@example.com");
        when(userRepository.findByEmail("TEST@EXAMPLE.COM")).thenReturn(Optional.of(user));
        when(jwtService.validateToken(request.getToken())).thenReturn(decodedJWT);

        AuthDTO.Response response = authService.validateToken(request);

        assertNotNull(response);
        assertEquals(UserDTO.from(user), response.getUser());
    }

    @Test
    void testValidateToken_whenUserNotFound_shouldThrowUserNotFoundException() {
        AuthDTO.TokenRequest request = AuthDTO.TokenRequest.builder()
                .token("validToken")
                .build();

        DecodedJWT decodedJWT = mock(DecodedJWT.class);
        when(decodedJWT.getSubject()).thenReturn("notfound@example.com");

        when(jwtService.validateToken(request.getToken())).thenReturn(decodedJWT);
        when(userRepository.findByEmail("NOTFOUND@EXAMPLE.COM")).thenReturn(Optional.empty());

        assertThrows(UserNotFoundException.class, () -> authService.validateToken(request));
    }
}

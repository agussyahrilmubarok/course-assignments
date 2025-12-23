package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.common.BaseServiceTest;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.DuplicateFieldException;
import com.example.witrack.backend.model.AuthDTO;
import com.example.witrack.backend.repos.UserRepository;
import com.example.witrack.backend.security.jwt.JwtProvider;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.util.Collections;
import java.util.Optional;
import java.util.UUID;

class AuthServiceImplTest extends BaseServiceTest {

    @InjectMocks
    private AuthServiceImpl authService;

    @Mock
    private UserRepository userRepository;

    @Mock
    private PasswordEncoder passwordEncoder;

    @Mock
    private JwtProvider jwtProvider;

    @Mock
    private AuthenticationManager authenticationManager;

    @Mock
    private Authentication authentication;

    private AuthDTO.SignUpRequest signUpRequest;
    private AuthDTO.SignInRequest signInRequest;
    private User testUser;

    @BeforeEach
    void setUp() {
        signUpRequest = AuthDTO.SignUpRequest.builder()
                .fullName("John Doe")
                .email("johndoe@test.com")
                .password("password123")
                .build();

        signInRequest = AuthDTO.SignInRequest.builder()
                .email("johndoe@test.com")
                .password("password123")
                .build();

        testUser = new User();
        testUser.setId(UUID.randomUUID());
        testUser.setFullName(signUpRequest.getFullName());
        testUser.setEmail(signUpRequest.getEmail().toLowerCase());
        testUser.setPassword("encodedPassword");
        testUser.setRoles(Collections.singleton(User.Role.ROLE_USER));
    }

    @Test
    void givenValidRequest_whenSignUp_thenReturnSignUpResponse() {
        Mockito.when(userRepository.existsByEmailIgnoreCase(Mockito.anyString())).thenReturn(false);
        Mockito.when(passwordEncoder.encode(Mockito.anyString())).thenReturn(testUser.getPassword());
        Mockito.when(userRepository.save(Mockito.any(User.class))).thenReturn(testUser);
        Mockito.when(jwtProvider.generateToken(Mockito.anyString(), Mockito.anySet())).thenReturn("jwt-token");

        AuthDTO.AuthResponse response = authService.signUp(signUpRequest);

        Assertions.assertNotNull(response);
        Assertions.assertEquals("jwt-token", response.getToken());
        Assertions.assertEquals(testUser.getEmail(), response.getUser().getEmail());
        Mockito.verify(userRepository, Mockito.times(1)).save(Mockito.any(User.class));
    }

    @Test
    void givenExistingEmail_whenSignUp_thenThrowDuplicateFieldException() {
        Mockito.when(userRepository.existsByEmailIgnoreCase(Mockito.anyString())).thenReturn(true);

        DuplicateFieldException exception = Assertions.assertThrows(DuplicateFieldException.class,
                () -> authService.signUp(signUpRequest));

        Assertions.assertEquals("Email is already in use", exception.getMessage());
        Assertions.assertEquals("email", exception.getField());
        Mockito.verify(userRepository, Mockito.never()).save(Mockito.any(User.class));
    }

    @Test
    void givenSaveFails_whenSignUp_thenThrowException() {
        Mockito.when(userRepository.existsByEmailIgnoreCase(Mockito.anyString())).thenReturn(false);
        Mockito.when(passwordEncoder.encode(Mockito.anyString())).thenReturn(testUser.getPassword());
        Mockito.when(userRepository.save(Mockito.any(User.class))).thenThrow(new RuntimeException("Database error"));

        RuntimeException exception = Assertions.assertThrows(RuntimeException.class,
                () -> authService.signUp(signUpRequest));

        Assertions.assertEquals("Database error", exception.getMessage());
        Mockito.verify(userRepository, Mockito.times(1)).save(Mockito.any(User.class));
    }

    @Test
    void givenValidRequest_whenSignIn_thenReturnSignInResponse() {
        Mockito.when(authenticationManager.authenticate(Mockito.any(UsernamePasswordAuthenticationToken.class)))
                .thenReturn(authentication);
        Mockito.when(jwtProvider.generateToken(authentication)).thenReturn("jwt-token");
        Mockito.when(userRepository.findByEmail(Mockito.anyString())).thenReturn(Optional.of(testUser));

        AuthDTO.AuthResponse response = authService.signIn(signInRequest);

        Assertions.assertNotNull(response);
        Assertions.assertEquals("jwt-token", response.getToken());
        Assertions.assertEquals(testUser.getEmail(), response.getUser().getEmail());
        Mockito.verify(userRepository, Mockito.times(1)).findByEmail(Mockito.anyString());
    }

    @Test
    void givenNonExistingEmail_whenSignIn_thenThrowUsernameNotFoundException() {
        Mockito.when(authenticationManager.authenticate(Mockito.any(UsernamePasswordAuthenticationToken.class)))
                .thenReturn(authentication);
        Mockito.when(jwtProvider.generateToken(authentication)).thenReturn("jwt-token");
        Mockito.when(userRepository.findByEmail(Mockito.anyString())).thenReturn(Optional.empty());

        UsernameNotFoundException exception = Assertions.assertThrows(UsernameNotFoundException.class,
                () -> authService.signIn(signInRequest));

        Assertions.assertTrue(exception.getMessage().contains("User is not found"));
        Mockito.verify(userRepository, Mockito.times(1)).findByEmail(Mockito.anyString());
    }

    @Test
    void givenInvalidCredentials_whenSignIn_thenThrowBadCredentialsException() {
        Mockito.when(authenticationManager.authenticate(Mockito.any(UsernamePasswordAuthenticationToken.class)))
                .thenThrow(new BadCredentialsException("Bad credentials"));

        Assertions.assertThrows(BadCredentialsException.class, () -> authService.signIn(signInRequest));
    }

    @Test
    void givenUnexpectedError_whenAuthenticate_thenThrowException() {
        Mockito.when(authenticationManager.authenticate(Mockito.any(UsernamePasswordAuthenticationToken.class)))
                .thenThrow(new RuntimeException("Unexpected error"));

        RuntimeException exception = Assertions.assertThrows(RuntimeException.class,
                () -> authService.signIn(signInRequest));

        Assertions.assertEquals("Unexpected error", exception.getMessage());
    }
}
package com.example.user.service.impl;

import com.example.user.domain.User;
import com.example.user.exception.UserNotFoundException;
import com.example.user.model.UserDTO;
import com.example.user.repos.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UserServiceImplTest {

    @InjectMocks
    private UserServiceImpl userService;

    @Mock
    private UserRepository userRepository;

    private User user;

    @BeforeEach
    void setUp() {
        user = new User();
        user.setId("123");
        user.setName("Test User");
        user.setEmail("test@example.com");
        user.setPassword("encodedPassword");
    }

    @Test
    void testFindById_success() {
        when(userRepository.findById("123")).thenReturn(Optional.of(user));

        UserDTO result = userService.findByID("123");

        assertNotNull(result);
        assertEquals("123", result.getId());
        assertEquals("Test User", result.getName());
        assertEquals("test@example.com", result.getEmail());

        verify(userRepository, times(1)).findById("123");
    }

    @Test
    void testFindById_userNotFound_shouldThrowException() {
        when(userRepository.findById("999")).thenReturn(Optional.empty());

        assertThrows(UserNotFoundException.class, () -> userService.findByID("999"));

        verify(userRepository, times(1)).findById("999");
    }
}

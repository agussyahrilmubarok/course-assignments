package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.model.UserResponse;
import com.example.witrack.backend.repository.UserRepository;
import com.example.witrack.backend.service.BaseServiceTest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;

import java.util.Collections;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

class UserServiceImplTest extends BaseServiceTest {

    @InjectMocks
    private UserServiceImpl userService;

    @Mock
    private UserRepository userRepository;

    private User testUser;

    @BeforeEach
    void setUp() {
        testUser = new User();
        testUser.setId(UUID.randomUUID().toString());
        testUser.setFullName("John Doe");
        testUser.setEmail("johndoe@mail.com");
        testUser.setPassword("encodedPassword");
        testUser.setRoles(Collections.singleton(User.Role.ROLE_USER));
    }

    @Test
    void givenExistingUserId_whenGetById_thenReturnUserResponse() {
        when(userRepository.findById(testUser.getId())).thenReturn(Optional.of(testUser));

        UserResponse response = userService.getById(testUser.getId());

        assertNotNull(response);
        assertEquals(testUser.getId(), response.getId());
        assertEquals(testUser.getFullName(), response.getFullName());
        assertEquals(testUser.getEmail(), response.getEmail());

        verify(userRepository, times(1)).findById(testUser.getId());
    }

    @Test
    void givenNonExistingUserId_whenGetById_thenThrowNotFoundException() {
        String nonExistingId = UUID.randomUUID().toString();
        when(userRepository.findById(nonExistingId)).thenReturn(Optional.empty());

        NotFoundException exception = assertThrows(NotFoundException.class,
                () -> userService.getById(nonExistingId));

        assertTrue(exception.getMessage().contains(nonExistingId));
        verify(userRepository, times(1)).findById(nonExistingId);
    }
}
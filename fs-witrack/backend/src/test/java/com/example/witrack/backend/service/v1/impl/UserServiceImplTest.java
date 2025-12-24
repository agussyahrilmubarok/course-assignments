package com.example.witrack.backend.service.v1.impl;

import com.example.witrack.backend.common.BaseServiceTest;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.model.UserDTO;
import com.example.witrack.backend.repos.UserRepository;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.Mockito;

import java.util.Collections;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertThrows;

class UserServiceImplTest extends BaseServiceTest {

    @InjectMocks
    private UserServiceImpl userService;

    @Mock
    private UserRepository userRepository;

    private User testUser;

    @BeforeEach
    void setUp() {
        testUser = new User();
        testUser.setId(UUID.randomUUID());
        testUser.setFullName("John Doe");
        testUser.setEmail("johndoe@mail.com");
        testUser.setPassword("encodedPassword");
        testUser.setRoles(Collections.singleton(User.Role.ROLE_USER));
    }

    @Test
    void givenExistingUserId_whenFindById_thenReturnUserResponse() {
        Mockito.when(userRepository.findById(Mockito.any())).thenReturn(Optional.of(testUser));

        UserDTO.UserResponse response = userService.findById(testUser.getId());

        Assertions.assertNotNull(response);
        Assertions.assertNotNull(response.getId());
        Assertions.assertEquals(testUser.getFullName(), response.getFullName());
        Assertions.assertEquals(testUser.getEmail(), response.getEmail());
        Mockito.verify(userRepository, Mockito.times(1)).findById(testUser.getId());
    }

    @Test
    void givenNonExistingUserId_whenFindById_thenThrowNotFoundException() {
        Mockito.when(userRepository.findById(Mockito.any())).thenReturn(Optional.empty());
        UUID nonExistingId = UUID.randomUUID();

        NotFoundException exception = assertThrows(NotFoundException.class,
                () -> userService.findById(nonExistingId));

        Assertions.assertTrue(exception.getMessage().contains(nonExistingId.toString()));
        Mockito.verify(userRepository, Mockito.times(1)).findById(nonExistingId);
    }
}
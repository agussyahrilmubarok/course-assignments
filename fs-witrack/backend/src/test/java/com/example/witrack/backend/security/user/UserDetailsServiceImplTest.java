package com.example.witrack.backend.security.user;

import com.example.witrack.backend.common.BaseServiceTest;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.repos.UserRepository;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;

import java.util.Collections;
import java.util.Optional;
import java.util.UUID;

class UserDetailsServiceImplTest extends BaseServiceTest {

    @InjectMocks
    private UserDetailsServiceImpl userDetailsService;

    @Mock
    private UserRepository userRepository;

    private User testUser;
    private UUID testUserId;

    @BeforeEach
    void setUp() {
        testUserId = UUID.randomUUID();
        testUser = new User();
        testUser.setId(testUserId);
        testUser.setEmail("johndoe@test.com");
        testUser.setPassword("encodedPassword");
        testUser.setRoles(Collections.singleton(User.Role.ROLE_USER));
    }

    @Test
    void givenExistingUserId_whenLoadUserByUsername_thenReturnUserDetails() {
        Mockito.when(userRepository.findById(testUserId)).thenReturn(Optional.of(testUser));

        UserDetails userDetails = userDetailsService.loadUserByUsername(testUserId.toString());

        Assertions.assertNotNull(userDetails);
        Assertions.assertEquals(testUser.getId(), ((UserDetailsImpl) userDetails).getId());
        Assertions.assertEquals(testUser.getEmail(), userDetails.getUsername());
        Assertions.assertEquals(testUser.getPassword(), userDetails.getPassword());
        Assertions.assertEquals(testUser.getRoles().size(), userDetails.getAuthorities().size());
        Mockito.verify(userRepository, Mockito.times(1)).findById(testUserId);
    }

    @Test
    void givenNonExistingUserId_whenLoadUserByUsername_thenThrowUsernameNotFoundException() {
        Mockito.when(userRepository.findById(testUserId)).thenReturn(Optional.empty());

        UsernameNotFoundException exception = Assertions.assertThrows(UsernameNotFoundException.class,
                () -> userDetailsService.loadUserByUsername(testUserId.toString()));

        Assertions.assertEquals("User not found", exception.getMessage());
        Mockito.verify(userRepository, Mockito.times(1)).findById(testUserId);
    }

    @Test
    void givenInvalidUUID_whenLoadUserByUsername_thenThrowIllegalArgumentException() {
        String invalidId = "invalid-uuid";

        Assertions.assertThrows(IllegalArgumentException.class,
                () -> userDetailsService.loadUserByUsername(invalidId));

        Mockito.verifyNoInteractions(userRepository);
    }
}
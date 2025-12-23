package com.example.witrack.backend.security;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.repository.UserRepository;
import com.example.witrack.backend.service.BaseServiceTest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;

import java.util.Collections;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

class UserDetailsServiceImplTest extends BaseServiceTest {

    @InjectMocks
    private UserDetailsServiceImpl userDetailsService;

    @Mock
    private UserRepository userRepository;

    private User mockUser;

    @BeforeEach
    void setUp() {
        mockUser = new User();
        mockUser.setId(UUID.randomUUID().toString());
        mockUser.setEmail("johndoe@mail.com");
        mockUser.setPassword("encodedPassword");
        mockUser.setRoles(Collections.singleton(User.Role.ROLE_USER));
    }

    @Test
    void loadUserByUsername_WhenUserExists_ShouldReturnUserDetails() {
        when(userRepository.findByEmail("johndoe@mail.com")).thenReturn(Optional.of(mockUser));

        UserDetails userDetails = userDetailsService.loadUserByUsername("johndoe@mail.com");

        assertNotNull(userDetails);
        assertEquals(mockUser.getEmail(), userDetails.getUsername());
        assertEquals(mockUser.getPassword(), userDetails.getPassword());
        assertTrue(userDetails.getAuthorities().stream()
                .anyMatch(a -> a.getAuthority().equals("ROLE_USER")));

        verify(userRepository, times(1)).findByEmail("johndoe@mail.com");
    }

    @Test
    void loadUserByUsername_WhenUserNotExists_ShouldThrowUsernameNotFoundException() {
        when(userRepository.findByEmail("unknown@mail.com")).thenReturn(Optional.empty());

        UsernameNotFoundException ex = assertThrows(UsernameNotFoundException.class,
                () -> userDetailsService.loadUserByUsername("unknown@mail.com"));

        assertEquals("User is not found", ex.getMessage());
        verify(userRepository, times(1)).findByEmail("unknown@mail.com");
    }
}
package com.example.witrack.backend.common;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.security.jwt.JwtProvider;
import com.example.witrack.backend.security.user.UserDetailsImpl;
import com.example.witrack.backend.security.user.UserDetailsServiceImpl;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc;
import org.springframework.context.annotation.Import;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashSet;
import java.util.UUID;

@Import(TestConfig.class)
@SpringBootTest
@AutoConfigureMockMvc
public abstract class BaseResourceTest extends TestcontainersConfig {

    @Autowired
    protected MockMvc mockMvc;

    @Autowired
    protected ObjectMapper objectMapper;

    @Autowired
    protected JwtProvider jwtProvider;

    @MockitoBean
    protected UserDetailsServiceImpl userDetailsService;

    protected User mockUser;
    protected User mockAdmin;

    @BeforeEach
    protected void initializeAuth() {
        this.mockUser = generateUser();
        final UserDetailsImpl mockUserDetails = new UserDetailsImpl(
                mockUser.getId(), mockUser.getEmail(), mockUser.getPassword(), mockUser.getRoles()
        );
        Mockito.when(userDetailsService.loadUserByUsername(mockUser.getId().toString()))
                .thenReturn(mockUserDetails);

        this.mockAdmin = generateAdmin();
        final UserDetailsImpl mockAdminDetails = new UserDetailsImpl(
                mockAdmin.getId(), mockAdmin.getEmail(), mockAdmin.getPassword(), mockAdmin.getRoles()
        );
        Mockito.when(userDetailsService.loadUserByUsername(mockAdmin.getId().toString()))
                .thenReturn(mockAdminDetails);
    }

    private User generateUser() {
        User user = new User();
        user.setId(UUID.randomUUID());
        user.setFullName("John Doe");
        user.setEmail("johndoe@mail.com");
        user.setPassword("encodedPassword");
        user.setRoles(Collections.singleton(User.Role.ROLE_USER));
        return user;
    }

    private User generateAdmin() {
        User user = new User();
        user.setId(UUID.randomUUID());
        user.setFullName("Jane Smith");
        user.setEmail("janesmith@mail.com");
        user.setPassword("encodedPassword");
        user.setRoles(new HashSet<>(Arrays.asList(User.Role.ROLE_USER, User.Role.ROLE_ADMIN)));
        return user;
    }

//    private String generateMockToken(UserDetailsImpl userDetails) {
//        return "Bearer " + jwtProvider.generateToken(userDetails.getId().toString(), userDetails.getRoles());
//    }
}

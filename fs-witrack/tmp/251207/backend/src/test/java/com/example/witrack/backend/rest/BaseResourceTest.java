package com.example.witrack.backend.rest;

import com.example.witrack.backend.TestContainerConfig;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.security.UserDetailsImpl;
import com.example.witrack.backend.security.UserDetailsServiceImpl;
import com.example.witrack.backend.security.jwt.JwtProvider;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashSet;
import java.util.UUID;

@SpringBootTest
@AutoConfigureMockMvc
public abstract class BaseResourceTest extends TestContainerConfig {

    @Autowired
    protected MockMvc mockMvc;

    @Autowired
    protected ObjectMapper objectMapper;

    @Autowired
    protected JwtProvider jwtProvider;

    @MockitoBean
    protected UserDetailsServiceImpl userDetailsService;

    protected User mockUser;
    //protected String mockUserToken;
    protected User mockAdmin;
    //protected String mockAdminToken;

    @BeforeEach
    protected void initializeAuth() {
        this.mockUser = generateUser();
        final UserDetailsImpl mockUserDetails = new UserDetailsImpl(
                mockUser.getId(), mockUser.getEmail(), mockUser.getPassword(), mockUser.getRoles()
        );
        //this.mockUserToken = generateMockToken(mockUserDetails);

        this.mockAdmin = generateAdmin();
        final UserDetailsImpl mockAdminDetails = new UserDetailsImpl(
                mockAdmin.getId(), mockAdmin.getEmail(), mockAdmin.getPassword(), mockAdmin.getRoles()
        );
        //this.mockAdminToken = generateMockToken(mockAdminDetails);

        Mockito.when(userDetailsService.loadUserByUsername(mockUser.getEmail()))
                .thenReturn(mockUserDetails);
        Mockito.when(userDetailsService.loadUserByUsername(mockAdmin.getEmail()))
                .thenReturn(mockAdminDetails);
    }

    private User generateUser() {
        User user = new User();
        user.setId(UUID.randomUUID().toString());
        user.setFullName("John Doe");
        user.setEmail("johndoe@mail.com");
        user.setPassword("encodedPassword");
        user.setRoles(Collections.singleton(User.Role.ROLE_USER));
        return user;
    }

    private User generateAdmin() {
        User user = new User();
        user.setId(UUID.randomUUID().toString());
        user.setFullName("Jane Smith");
        user.setEmail("janesmith@mail.com");
        user.setPassword("encodedPassword");
        user.setRoles(new HashSet<>(Arrays.asList(User.Role.ROLE_USER, User.Role.ROLE_ADMIN)));
        return user;
    }

//    private String generateMockToken(UserDetailsImpl userDetails) {
//        String id = userDetails.getId();
//        List<String> roles = userDetails.getAuthorities()
//                .stream()
//                .map(auth -> auth.getAuthority())
//                .collect(Collectors.toList());
//
//        Map<String, Object> claims = new HashMap<>();
//        claims.put("id", id);
//        claims.put("roles", roles);
//
//        return "Bearer " + jwtProvider.generateToken(claims, userDetails.getUsername());
//    }
}

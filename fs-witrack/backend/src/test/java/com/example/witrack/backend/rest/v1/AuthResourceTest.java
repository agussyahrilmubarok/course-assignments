package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.model.*;
import com.example.witrack.backend.rest.BaseResourceTest;
import com.example.witrack.backend.service.impl.AuthServiceImpl;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;

import java.util.Collections;
import java.util.UUID;

class AuthResourceTest extends BaseResourceTest {

    @MockitoBean
    private AuthServiceImpl authService;

    @Test
    void givenSignUpRequest_ReturnSignUpResponse() throws Exception {
        SignUpRequest request = SignUpRequest.builder()
                .fullName("Chunk Smith")
                .email("chunksmith@mail.com")
                .password("secretpassword")
                .build();

        User dummyUser = new User();
        dummyUser.setId(UUID.randomUUID().toString());
        dummyUser.setFullName(request.getFullName());
        dummyUser.setEmail(request.getEmail());
        dummyUser.setRoles(Collections.singleton(User.Role.ROLE_USER));

        SignUpResponse response = SignUpResponse.builder()
                .token("jwtToken")
                .user(UserResponse.fromUser(dummyUser))
                .build();

        Mockito.when(authService.signUp(request)).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/auth/sign-up")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.token").exists())
                .andExpect(MockMvcResultMatchers.jsonPath("$.user.fullName").value("Chunk Smith"))
                .andExpect(MockMvcResultMatchers.jsonPath("$.user.email").value("chunksmith@mail.com"));

        Mockito.verify(authService).signUp(request);
    }

    @Test
    void givenSignInRequest_ReturnSignInResponse() throws Exception {
        SignInRequest request = SignInRequest.builder()
                .email("chunksmith@mail.com")
                .password("secretpassword")
                .build();

        User dummyUser = new User();
        dummyUser.setId(UUID.randomUUID().toString());
        dummyUser.setFullName("Chunk Smith");
        dummyUser.setEmail(request.getEmail());
        dummyUser.setRoles(Collections.singleton(User.Role.ROLE_USER));

        SignInResponse response = SignInResponse.builder()
                .token("jwtToken")
                .user(UserResponse.fromUser(dummyUser))
                .build();

        Mockito.when(authService.signIn(request)).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/auth/sign-in")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.token").exists())
                .andExpect(MockMvcResultMatchers.jsonPath("$.user.fullName").value("Chunk Smith"))
                .andExpect(MockMvcResultMatchers.jsonPath("$.user.email").value("chunksmith@mail.com"));

        Mockito.verify(authService).signIn(request);
    }
}
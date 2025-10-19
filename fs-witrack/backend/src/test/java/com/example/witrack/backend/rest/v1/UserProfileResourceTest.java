package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.UserResponse;
import com.example.witrack.backend.rest.BaseResourceTest;
import com.example.witrack.backend.security.CurrentUserDetails;
import com.example.witrack.backend.service.impl.UserServiceImpl;
import lombok.SneakyThrows;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.http.MediaType;
import org.springframework.security.test.context.support.WithMockUser;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;

class UserProfileResourceTest extends BaseResourceTest {

    @MockitoBean
    private UserServiceImpl userService;

    @MockitoBean
    private CurrentUserDetails currentUserDetails;

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenGetMeRequest_ReturnUserResponse() {
        UserResponse mockResponse = UserResponse.fromUser(mockUser);

        Mockito.when(currentUserDetails.getId()).thenReturn(mockUser.getId());
        Mockito.when(userService.getById(mockUser.getId())).thenReturn(mockResponse);

        mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/users/profiles/me")
                        .contentType(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.id").value(mockUser.getId()))
                .andExpect(MockMvcResultMatchers.jsonPath("$.fullName").value(mockUser.getFullName()))
                .andExpect(MockMvcResultMatchers.jsonPath("$.email").value(mockUser.getEmail()));

        Mockito.verify(userService, Mockito.times(1)).getById(mockUser.getId());
    }
}
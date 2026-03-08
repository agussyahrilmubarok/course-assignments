package com.example.witrack.backend.model;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class SignInResponse {

    private String token;

    private UserResponse user;
}

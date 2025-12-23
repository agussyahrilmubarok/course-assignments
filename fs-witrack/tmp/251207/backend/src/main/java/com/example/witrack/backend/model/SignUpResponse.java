package com.example.witrack.backend.model;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class SignUpResponse {

    private String token;

    private UserResponse user;
}

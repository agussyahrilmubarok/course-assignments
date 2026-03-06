package com.example.witrack.backend.service.v1;

import com.example.witrack.backend.model.AuthDTO;

public interface AuthService {

    AuthDTO.AuthResponse signUp(AuthDTO.SignUpRequest param);

    AuthDTO.AuthResponse signIn(AuthDTO.SignInRequest param);
}

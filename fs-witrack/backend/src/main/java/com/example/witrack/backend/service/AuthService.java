package com.example.witrack.backend.service;

import com.example.witrack.backend.model.SignInRequest;
import com.example.witrack.backend.model.SignInResponse;
import com.example.witrack.backend.model.SignUpRequest;
import com.example.witrack.backend.model.SignUpResponse;

public interface AuthService {

    SignUpResponse signUp(SignUpRequest request);

    SignInResponse signIn(SignInRequest request);
}

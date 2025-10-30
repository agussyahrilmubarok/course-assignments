package com.example.user.service;

import com.example.user.model.AuthDTO;

public interface AuthService {

    void signUp(final AuthDTO.SignUp param);

    AuthDTO.Response signIn(final AuthDTO.SignIn param);

    AuthDTO.Response validateToken(final AuthDTO.TokenRequest param);
}

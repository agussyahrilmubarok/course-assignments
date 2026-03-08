package com.example.user.service;

import com.auth0.jwt.interfaces.DecodedJWT;
import com.example.user.domain.User;

public interface JWTService {

    String generate(User user);

    DecodedJWT validateToken(String token);
}

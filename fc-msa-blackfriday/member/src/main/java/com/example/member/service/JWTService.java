package com.example.member.service;

import com.auth0.jwt.interfaces.DecodedJWT;
import com.example.member.domain.User;

public interface JWTService {

    String generate(User user);

    DecodedJWT validateToken(String token);
}

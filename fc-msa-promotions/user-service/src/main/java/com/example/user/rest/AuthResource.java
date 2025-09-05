package com.example.user.rest;

import com.example.user.service.AuthService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/api/v1/auth", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class AuthResource {

    private AuthService authService;

    @PostMapping("/sign-up")
    public ResponseEntity<Void> signUp() {
        return new ResponseEntity<>(null);
    }

    @PostMapping("/sign-in")
    public ResponseEntity<Void> signIn() {
        return new ResponseEntity<>(null);
    }
}

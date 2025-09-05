package com.example.user.rest;

import com.example.user.model.AuthDTO;
import com.example.user.service.AuthService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/api/v1/auth", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class AuthResource {

    private final AuthService authService;

    @PostMapping("/sign-up")
    public ResponseEntity<Void> signUp(@RequestBody @Valid final AuthDTO.SignUp payload) {
        authService.signUp(payload);
        return ResponseEntity.status(HttpStatus.CREATED).build();
    }

    @PostMapping("/sign-in")
    public ResponseEntity<AuthDTO.Response> signIn(@RequestBody @Valid final AuthDTO.SignIn payload) {
        return ResponseEntity.ok(authService.signIn(payload));
    }

    @PostMapping("/validate-token")
    public ResponseEntity<AuthDTO.Response> validateToken(@RequestBody @Valid final AuthDTO.TokenRequest payload) {
        return ResponseEntity.ok(authService.validateToken(payload));
    }
}

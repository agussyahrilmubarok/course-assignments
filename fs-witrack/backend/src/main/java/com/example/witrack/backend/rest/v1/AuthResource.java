package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.AuthDTO;
import com.example.witrack.backend.service.v1.AuthService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController("AuthResourceV1")
@RequestMapping(value = "/api/v1/auth", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class AuthResource {

    private final AuthService authService;

    @PostMapping("/sign-up")
    public ResponseEntity<AuthDTO.AuthResponse> signUp(@Valid @RequestBody AuthDTO.SignUpRequest request) {
        return ResponseEntity.ok(authService.signUp(request));
    }

    @PostMapping("/sign-in")
    public ResponseEntity<AuthDTO.AuthResponse> signIn(@Valid @RequestBody AuthDTO.SignInRequest request) {
        return ResponseEntity.ok(authService.signIn(request));
    }
}

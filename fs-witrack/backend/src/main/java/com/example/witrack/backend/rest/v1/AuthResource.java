package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.SignInRequest;
import com.example.witrack.backend.model.SignInResponse;
import com.example.witrack.backend.model.SignUpRequest;
import com.example.witrack.backend.model.SignUpResponse;
import com.example.witrack.backend.service.AuthService;
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
    public ResponseEntity<SignUpResponse> signUp(@Valid @RequestBody SignUpRequest request) {
        SignUpResponse response = authService.signUp(request);
        return ResponseEntity.ok(response);
    }

    @PostMapping("/sign-in")
    public ResponseEntity<SignInResponse> signIn(@Valid @RequestBody SignInRequest request) {
        SignInResponse response = authService.signIn(request);
        return ResponseEntity.ok(response);
    }
}

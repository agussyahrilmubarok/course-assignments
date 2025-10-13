package com.example.member.rest;

import com.example.member.model.MemberDTO;
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
@RequestMapping(value = "/api/v1/members", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class MemberResource {

    @PostMapping("/sign-up")
    public ResponseEntity<Void> signUp(@RequestBody @Valid final MemberDTO.SignUp payload) {
        return ResponseEntity.status(HttpStatus.CREATED).build();
    }

    @PostMapping("/sign-in")
    public ResponseEntity<MemberDTO.ResponseWithToken> signUp(@RequestBody @Valid final MemberDTO.SignIn payload) {
        return ResponseEntity.status(HttpStatus.CREATED).build();
    }

    @PostMapping("/validate")
    public ResponseEntity<MemberDTO.Response> validate(@RequestBody @Valid final MemberDTO.SignIn payload) {
        return ResponseEntity.status(HttpStatus.CREATED).build();
    }
}

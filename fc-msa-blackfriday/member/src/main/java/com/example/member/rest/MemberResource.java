package com.example.member.rest;

import com.example.member.model.MemberDTO;
import com.example.member.service.MemberService;
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

    private final MemberService memberService;

    @PostMapping("/sign-up")
    public ResponseEntity<Void> signUp(@RequestBody @Valid final MemberDTO.SignUp payload) {
        memberService.signUp(payload);
        return ResponseEntity.status(HttpStatus.CREATED).build();
    }

    @PostMapping("/sign-in")
    public ResponseEntity<MemberDTO.ResponseWithToken> signUp(@RequestBody @Valid final MemberDTO.SignIn payload) {
        return ResponseEntity.ok(memberService.signIn(payload));
    }

    @PostMapping("/validate")
    public ResponseEntity<MemberDTO.ResponseWithToken> validate(@RequestBody @Valid final MemberDTO.ValidateToken payload) {
        return ResponseEntity.ok(memberService.validateToken(payload.getToken()));
    }
}

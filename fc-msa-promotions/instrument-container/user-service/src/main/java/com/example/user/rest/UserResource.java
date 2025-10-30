package com.example.user.rest;

import com.example.user.model.UserDTO;
import com.example.user.service.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/api/v1/users", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class UserResource {

    private final UserService userService;

    @GetMapping("/me")
    public ResponseEntity<UserDTO> getMe(@RequestHeader("X-USER-ID") String userId) {
        return new ResponseEntity<>(userService.findByID(userId), HttpStatus.OK);
    }
}

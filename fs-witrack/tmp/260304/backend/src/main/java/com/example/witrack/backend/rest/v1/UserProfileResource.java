package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.UserDTO;
import com.example.witrack.backend.security.user.CurrentUserDetails;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController("UserProfileResourceV1")
@RequestMapping(value = "/api/v1/users/profiles", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class UserProfileResource {

    private final CurrentUserDetails currentUserDetails;

    @GetMapping("/me")
    public ResponseEntity<UserDTO.UserResponse> getMe() {
        return ResponseEntity.ok(UserDTO.UserResponse.fromUser(currentUserDetails.getUser()));
    }
}

package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.UserResponse;
import com.example.witrack.backend.security.CurrentUserDetails;
import com.example.witrack.backend.service.UserService;
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

    private final UserService userService;
    private final CurrentUserDetails currentUserDetails;

    @GetMapping("/me")
    public ResponseEntity<UserResponse> getMe() {
        String id = currentUserDetails.getId();
        UserResponse response = userService.getById(id);
        return ResponseEntity.ok(response);
    }
}

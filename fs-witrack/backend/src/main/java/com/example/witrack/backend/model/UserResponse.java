package com.example.witrack.backend.model;

import com.example.witrack.backend.domain.User;
import lombok.Builder;
import lombok.Data;

import java.util.List;
import java.util.stream.Collectors;

@Data
@Builder
public class UserResponse {

    private String id;

    private String fullName;

    private String email;

    private List<String> roles;

    public static UserResponse fromUser(User user) {
        if (user == null) {
            return null;
        }

        return UserResponse.builder()
                .id(user.getId())
                .fullName(user.getFullName())
                .email(user.getEmail())
                .roles(user.getRoles() != null
                        ? user.getRoles().stream()
                        .map(Enum::name)
                        .collect(Collectors.toList())
                        : null)
                .build();
    }
}

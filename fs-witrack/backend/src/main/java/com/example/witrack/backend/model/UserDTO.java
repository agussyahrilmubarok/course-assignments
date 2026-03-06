package com.example.witrack.backend.model;

import com.example.witrack.backend.domain.User;
import com.fasterxml.jackson.annotation.JsonInclude;
import lombok.*;

import java.util.List;
import java.util.stream.Collectors;

public class UserDTO {

    @Data
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class UserResponse {

        private String id;
        private String fullName;
        private String email;
        private List<String> roles;

        public static UserDTO.UserResponse fromUser(User user) {
            if (user == null) {
                return null;
            }

            return UserDTO.UserResponse.builder()
                    .id(user.getId().toString())
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

}

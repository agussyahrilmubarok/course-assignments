package com.example.witrack.backend.domain;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.Id;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.annotation.Version;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.OffsetDateTime;
import java.util.Set;

@Document("users")
@Getter
@Setter
public class User {

    @Id
    @NotNull
    private String id;

    @NotNull
    @Size(max = 255)
    private String fullName;

    @Indexed(unique = true)
    @NotNull
    @Size(max = 255)
    private String email;

    @NotNull
    @Size(max = 255)
    private String password;

    @NotNull
    private Set<Role> roles;

    @CreatedDate
    private OffsetDateTime createdAt;

    @LastModifiedDate
    private OffsetDateTime updatedAt;

    @Version
    private Integer version;

    public boolean hasRole(Role role) {
        return roles != null && roles.contains(role);
    }

    public enum Role {
        ROLE_USER,
        ROLE_ADMIN
    }
}

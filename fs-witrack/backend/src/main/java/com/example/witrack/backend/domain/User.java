package com.example.witrack.backend.domain;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.OffsetDateTime;
import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

@Entity
@Table(name = "users", indexes = {@Index(name = "uk_user_email", columnList = "email", unique = true)})
@Getter
@Setter
public class User {

    @Id
    @GeneratedValue
    private UUID id;

    @NotNull
    @Size(max = 255)
    @Column(nullable = false, length = 255)
    private String fullName;

    @NotNull
    @Size(max = 255)
    @Column(nullable = false, unique = true, length = 255)
    private String email;

    @NotNull
    @Size(max = 255)
    @Column(nullable = false, length = 255)
    private String password;

    @ElementCollection(fetch = FetchType.EAGER)
    @CollectionTable(
            name = "user_roles",
            joinColumns = @JoinColumn(name = "user_id")
    )
    @Enumerated(EnumType.STRING)
    @Column(name = "role", nullable = false, length = 50)
    private Set<Role> roles = new HashSet<>();

    @CreationTimestamp
    @Column(updatable = false)
    private OffsetDateTime createdAt;

    @UpdateTimestamp
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

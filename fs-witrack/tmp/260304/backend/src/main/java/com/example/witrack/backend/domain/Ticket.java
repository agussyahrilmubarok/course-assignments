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
@Table(name = "tickets", indexes = {@Index(name = "uk_ticket_code", columnList = "code", unique = true)})
@Getter
@Setter
public class Ticket {

    @Id
    @GeneratedValue
    private UUID id;

    @NotNull
    @Size(max = 255)
    @Column(nullable = false, unique = true, length = 255)
    private String code;

    @NotNull
    @Size(max = 255)
    @Column(nullable = false, length = 255)
    private String title;

    @NotNull
    @Column(nullable = false, columnDefinition = "TEXT")
    private String description;

    @NotNull
    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 50)
    private Status status;

    @NotNull
    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 50)
    private Priority priority;

    private OffsetDateTime completeAt;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    @OneToMany(
            mappedBy = "ticket",
            fetch = FetchType.LAZY,
            cascade = CascadeType.ALL,
            orphanRemoval = true
    )
    private Set<TicketComment> ticketComments = new HashSet<>();

    @CreationTimestamp
    @Column(updatable = false)
    private OffsetDateTime createdAt;

    @UpdateTimestamp
    private OffsetDateTime updatedAt;

    @Version
    private Integer version;

    public enum Status {
        OPEN,
        IN_PROGRESS,
        RESOLVED,
        REJECTED
    }

    public enum Priority {
        LOW,
        MEDIUM,
        HIGH,
        CRITICAL
    }
}

package com.example.witrack.backend.domain;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.*;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.DocumentReference;

import java.time.OffsetDateTime;
import java.util.HashSet;
import java.util.Set;

@Document("tickets")
@Getter
@Setter
public class Ticket {

    @Id
    @NotNull
    private String id;

    @Indexed(unique = true)
    @NotNull
    @Size(max = 255)
    private String code;

    @NotNull
    @Size(max = 255)
    private String title;

    @NotNull
    private String description;

    @NotNull
    private Status status;

    @NotNull
    private Priority priority;

    private OffsetDateTime completeAt;

    @DocumentReference(lazy = true)
    @NotNull
    private User user;

    @DocumentReference(lazy = true, lookup = "{ 'ticket' : ?#{#self._id} }")
    @ReadOnlyProperty
    private Set<TicketReply> ticketReplies = new HashSet<>();

    @CreatedDate
    private OffsetDateTime createdAt;

    @LastModifiedDate
    private OffsetDateTime updatedAt;

    @Version
    private Integer version;

    public enum Status {
        OPEN,
        ONPROGRESS,
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

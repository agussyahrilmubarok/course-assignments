package com.example.witrack.backend.model;

import com.example.witrack.backend.domain.Ticket;
import com.fasterxml.jackson.annotation.JsonInclude;
import lombok.Builder;
import lombok.Data;

import java.time.OffsetDateTime;

@Data
@Builder
@JsonInclude(JsonInclude.Include.NON_NULL)
public class TicketResponse {

    private String id;

    private String code;

    private String title;

    private String description;

    private String status;

    private String priority;

    private OffsetDateTime completedAt;

    private OffsetDateTime createdAt;

    private OffsetDateTime updatedAt;

    private UserResponse user;

    private Long totalReplies;

    public static TicketResponse fromTicket(Ticket ticket) {
        if (ticket == null) {
            return null;
        }

        return TicketResponse.builder()
                .id(ticket.getId())
                .code(ticket.getCode())
                .title(ticket.getTitle())
                .description(ticket.getDescription() != null ? ticket.getDescription() : null)
                .status(ticket.getStatus() != null ? ticket.getStatus().name() : null)
                .priority(ticket.getPriority() != null ? ticket.getPriority().name() : null)
                .completedAt(ticket.getCompleteAt() != null ? ticket.getCompleteAt() : null)
                .createdAt(ticket.getCreatedAt())
                .updatedAt(ticket.getUpdatedAt())
                .user(ticket.getUser() != null ? UserResponse.fromUser(ticket.getUser()) : null)
                .totalReplies(ticket.getTicketReplies() != null ? ticket.getTicketReplies().stream().count() : 0)
                .build();
    }
}

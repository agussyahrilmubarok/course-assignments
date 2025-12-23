package com.example.witrack.backend.model;

import com.example.witrack.backend.domain.Ticket;
import com.fasterxml.jackson.annotation.JsonInclude;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.*;

import java.time.OffsetDateTime;

@Getter
@Setter
public class TicketDTO {

    @Data
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class TicketRequest {
        @NotBlank(message = "Title is required")
        @Size(max = 100, message = "Title must not exceed 100 characters")
        private String title;

        @Size(max = 500, message = "Description must not exceed 500 characters")
        private String description;

        @NotBlank(message = "Status is required")
        private String status;

        @NotBlank(message = "Priority is required")
        private String priority;
    }

    @Data
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class TicketResponse {
        private String id;
        private String code;
        private String title;
        private String description;
        private String status;
        private String priority;
        private OffsetDateTime completedAt;
        private OffsetDateTime createdAt;
        private OffsetDateTime updatedAt;
        private UserDTO.UserResponse user;
        private Long totalReplies;

        public static TicketResponse fromTicket(Ticket ticket) {
            if (ticket == null) {
                return null;
            }

            return TicketResponse.builder()
                    .id(ticket.getId().toString())
                    .code(ticket.getCode())
                    .title(ticket.getTitle())
                    .description(ticket.getDescription() != null ? ticket.getDescription() : null)
                    .status(ticket.getStatus() != null ? ticket.getStatus().name() : null)
                    .priority(ticket.getPriority() != null ? ticket.getPriority().name() : null)
                    .completedAt(ticket.getCompleteAt() != null ? ticket.getCompleteAt() : null)
                    .createdAt(ticket.getCreatedAt())
                    .updatedAt(ticket.getUpdatedAt())
                    .user(ticket.getUser() != null ? UserDTO.UserResponse.fromUser(ticket.getUser()) : null)
                    .totalReplies(ticket.getTicketComments() != null ? (long) ticket.getTicketComments().size() : 0)
                    .build();
        }
    }
}

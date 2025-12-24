package com.example.witrack.backend.model;

import com.example.witrack.backend.domain.TicketComment;
import com.fasterxml.jackson.annotation.JsonInclude;
import jakarta.validation.constraints.Size;
import lombok.*;

import java.time.OffsetDateTime;

@Getter
@Setter
public class TicketCommentDTO {

    @Data
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class TicketCommentRequest {
        @Size(min = 10, max = 500, message = "Content must be between 10 and 500 characters")
        private String content;
    }

    @Data
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class TicketCommentResponse {
        private String id;
        private String content;
        private OffsetDateTime createdAt;
        private OffsetDateTime updatedAt;
        private UserDTO.UserResponse user;

        public static TicketCommentResponse fromTicketComment(TicketComment ticketComment) {
            if (ticketComment == null) {
                return null;
            }

            return TicketCommentResponse.builder()
                    .id(ticketComment.getId().toString())
                    .content(ticketComment.getContent())
                    .createdAt(ticketComment.getCreatedAt())
                    .updatedAt(ticketComment.getUpdatedAt())
                    .user(ticketComment.getUser() != null ? UserDTO.UserResponse.fromUser(ticketComment.getUser()) : null)
                    .build();
        }
    }
}

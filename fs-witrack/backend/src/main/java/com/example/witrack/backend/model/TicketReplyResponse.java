package com.example.witrack.backend.model;

import com.example.witrack.backend.domain.TicketReply;
import lombok.Builder;
import lombok.Data;

import java.time.OffsetDateTime;

@Data
@Builder
public class TicketReplyResponse {

    private String id;

    private String content;

    private OffsetDateTime createdAt;

    private OffsetDateTime updatedAt;

    private UserResponse user;

    public static TicketReplyResponse fromTicketReply(TicketReply ticketReply) {
        if (ticketReply == null) {
            return null;
        }

        return TicketReplyResponse.builder()
                .id(ticketReply.getId())
                .content(ticketReply.getContent())
                .createdAt(ticketReply.getCreatedAt())
                .updatedAt(ticketReply.getUpdatedAt())
                .user(ticketReply.getUser() != null ? UserResponse.fromUser(ticketReply.getUser()) : null)
                .build();
    }
}

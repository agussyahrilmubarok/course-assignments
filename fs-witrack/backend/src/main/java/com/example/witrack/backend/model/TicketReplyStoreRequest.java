package com.example.witrack.backend.model;

import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class TicketReplyStoreRequest {

    @Size(max = 500, message = "Content must not exceed 500 characters")
    private String content;
}

package com.example.witrack.backend.model;

import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class TicketReplyStoreRequest {

    @Size(min = 10, max = 500, message = "Content must be between 10 and 500 characters")
    private String content;
}

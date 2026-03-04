package com.example.witrack.backend.service.v1;

import com.example.witrack.backend.model.TicketCommentDTO;

public interface TicketCommentService {

    TicketCommentDTO.TicketCommentResponse create(String code, TicketCommentDTO.TicketCommentRequest request);
}

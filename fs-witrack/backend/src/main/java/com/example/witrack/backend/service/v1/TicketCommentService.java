package com.example.witrack.backend.service;

import com.example.witrack.backend.model.TicketCommentDTO;

public interface TicketCommentService {

    TicketCommentDTO.TicketCommentResponse create(String code, TicketCommentDTO.TicketCommentRequest request);
}

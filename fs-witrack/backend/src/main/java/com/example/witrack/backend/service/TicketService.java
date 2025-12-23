package com.example.witrack.backend.service;

import com.example.witrack.backend.model.TicketDTO;

import java.util.UUID;

public interface TicketService {

    TicketDTO.TicketResponse create(TicketDTO.TicketRequest request);

    TicketDTO.TicketResponse update(UUID id, TicketDTO.TicketRequest request);

    void delete(UUID id);
}

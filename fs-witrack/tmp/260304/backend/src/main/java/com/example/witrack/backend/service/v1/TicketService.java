package com.example.witrack.backend.service.v1;

import com.example.witrack.backend.model.TicketDTO;

import java.util.List;
import java.util.UUID;

public interface TicketService {

    TicketDTO.TicketResponse create(TicketDTO.TicketRequest request);

    TicketDTO.TicketResponse update(UUID id, TicketDTO.TicketRequest request);

    void delete(UUID id);

    TicketDTO.TicketResponse findById(UUID id);

    TicketDTO.TicketResponse findByCode(String code);

    List<TicketDTO.TicketResponse> searchTicket(String search, String status, String priority, String date);

    List<TicketDTO.TicketResponse> searchMyTicket(String search, String status, String priority, String date);
}

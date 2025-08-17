package com.example.witrack.backend.service;

import com.example.witrack.backend.model.*;

import java.util.List;

public interface TicketService {

    List<TicketResponse> getTickets(String search, String status, String priority);

    TicketDetailResponse getTicketByCode(String code);

    TicketResponse createTicket(TicketStoreRequest request);

    TicketResponse updateTicket(String code, TicketStoreRequest request);

    void deleteTicketByCode(String code);

    TicketReplyResponse createTicketReply(String code, TicketReplyStoreRequest request);
}

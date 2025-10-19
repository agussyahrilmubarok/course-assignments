package com.example.witrack.backend.service;

import com.example.witrack.backend.model.*;

import java.util.List;

public interface TicketService {

    List<TicketResponse> getTickets(String keyword, String status, String priority, String date);

    List<TicketResponse> getMyTickets(String keyword, String status, String priority, String date);

    TicketDetailResponse getTicketByCode(String code);

    TicketResponse createTicket(TicketStoreRequest request);

    TicketResponse updateTicket(String code, TicketStoreRequest request);

    void deleteTicketByCode(String code);

    TicketReplyResponse createTicketReply(String code, TicketReplyStoreRequest request);
}

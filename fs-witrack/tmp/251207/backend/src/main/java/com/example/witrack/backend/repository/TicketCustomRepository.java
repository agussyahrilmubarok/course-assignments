package com.example.witrack.backend.repository;

import com.example.witrack.backend.domain.Ticket;

import java.time.OffsetDateTime;
import java.util.List;

public interface TicketCustomRepository {

    List<Ticket> searchTickets(String keyword,
                               Ticket.Status status,
                               Ticket.Priority priority,
                               OffsetDateTime startAt,
                               OffsetDateTime endAt);

    List<Ticket> searchTickets(String keyword,
                               Ticket.Status status,
                               Ticket.Priority priority,
                               OffsetDateTime startAt,
                               OffsetDateTime endAt,
                               String userID);
}

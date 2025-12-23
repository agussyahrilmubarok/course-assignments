package com.example.witrack.backend.repos;

import com.example.witrack.backend.domain.Ticket;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface TicketRepository extends JpaRepository<Ticket, UUID> {
}

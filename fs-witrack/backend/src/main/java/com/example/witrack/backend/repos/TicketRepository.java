package com.example.witrack.backend.repos;

import com.example.witrack.backend.domain.Ticket;
import java.util.UUID;
import org.springframework.data.jpa.repository.JpaRepository;


public interface TicketRepository extends JpaRepository<Ticket, UUID> {

    Ticket findFirstByUserId(UUID id);

    boolean existsByCodeIgnoreCase(String code);

}

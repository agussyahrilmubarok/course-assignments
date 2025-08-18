package com.example.witrack.backend.repository;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.User;
import org.springframework.data.mongodb.repository.MongoRepository;

import java.util.List;
import java.util.Optional;


public interface TicketRepository extends MongoRepository<Ticket, String> {

    List<Ticket> findByStatus(String status);

    List<Ticket> findByPriority(String priority);

    List<Ticket> findByStatusAndPriority(String status, String priority);

    List<Ticket> findByCodeIgnoreCaseContainingOrTitleIgnoreCaseContainingOrDescriptionIgnoreCaseContaining(
            String code, String title, String description
    );

    Optional<Ticket> findByCode(String code);

    Ticket findFirstByUser(User user);

    boolean existsByIdIgnoreCase(String id);

    boolean existsByCodeIgnoreCase(String code);
}

package com.example.witrack.backend.repository;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.TicketReply;
import com.example.witrack.backend.domain.User;
import org.springframework.data.mongodb.repository.MongoRepository;


public interface TicketReplyRepository extends MongoRepository<TicketReply, String> {

    TicketReply findFirstByUser(User user);

    TicketReply findFirstByTicket(Ticket ticket);

    boolean existsByIdIgnoreCase(String id);
}

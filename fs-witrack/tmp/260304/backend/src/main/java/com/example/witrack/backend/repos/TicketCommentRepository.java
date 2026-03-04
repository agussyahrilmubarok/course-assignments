package com.example.witrack.backend.repos;

import com.example.witrack.backend.domain.TicketComment;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface TicketCommentRepository extends JpaRepository<TicketComment, UUID> {
}

package com.example.witrack.backend.repos;

import com.example.witrack.backend.domain.Comment;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface CommentRepository extends JpaRepository<Comment, UUID> {

    Comment findFirstByTicketId(UUID id);

    Comment findFirstByUserId(UUID id);

}

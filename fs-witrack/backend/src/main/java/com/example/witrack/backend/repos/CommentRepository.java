package com.example.witrack.backend.repos;

import com.example.witrack.backend.domain.Comment;
import java.util.UUID;
import org.springframework.data.jpa.repository.JpaRepository;


public interface CommentRepository extends JpaRepository<Comment, UUID> {

    Comment findFirstByTicketId(UUID id);

    Comment findFirstByUserId(UUID id);

}

package com.example.witrack.backend.service;

import com.example.witrack.backend.domain.Comment;
import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.events.BeforeDeleteTicket;
import com.example.witrack.backend.events.BeforeDeleteUser;
import com.example.witrack.backend.model.CommentDTO;
import com.example.witrack.backend.repos.CommentRepository;
import com.example.witrack.backend.repos.TicketRepository;
import com.example.witrack.backend.repos.UserRepository;
import com.example.witrack.backend.util.NotFoundException;
import com.example.witrack.backend.util.ReferencedException;
import java.util.List;
import java.util.UUID;
import org.springframework.context.event.EventListener;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;


@Service
public class CommentService {

    private final CommentRepository commentRepository;
    private final TicketRepository ticketRepository;
    private final UserRepository userRepository;

    public CommentService(final CommentRepository commentRepository,
            final TicketRepository ticketRepository, final UserRepository userRepository) {
        this.commentRepository = commentRepository;
        this.ticketRepository = ticketRepository;
        this.userRepository = userRepository;
    }

    public List<CommentDTO> findAll() {
        final List<Comment> comments = commentRepository.findAll(Sort.by("id"));
        return comments.stream()
                .map(comment -> mapToDTO(comment, new CommentDTO()))
                .toList();
    }

    public CommentDTO get(final UUID id) {
        return commentRepository.findById(id)
                .map(comment -> mapToDTO(comment, new CommentDTO()))
                .orElseThrow(NotFoundException::new);
    }

    public UUID create(final CommentDTO commentDTO) {
        final Comment comment = new Comment();
        mapToEntity(commentDTO, comment);
        return commentRepository.save(comment).getId();
    }

    public void update(final UUID id, final CommentDTO commentDTO) {
        final Comment comment = commentRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        mapToEntity(commentDTO, comment);
        commentRepository.save(comment);
    }

    public void delete(final UUID id) {
        final Comment comment = commentRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        commentRepository.delete(comment);
    }

    private CommentDTO mapToDTO(final Comment comment, final CommentDTO commentDTO) {
        commentDTO.setId(comment.getId());
        commentDTO.setContent(comment.getContent());
        commentDTO.setTicket(comment.getTicket() == null ? null : comment.getTicket().getId());
        commentDTO.setUser(comment.getUser() == null ? null : comment.getUser().getId());
        return commentDTO;
    }

    private Comment mapToEntity(final CommentDTO commentDTO, final Comment comment) {
        comment.setContent(commentDTO.getContent());
        final Ticket ticket = commentDTO.getTicket() == null ? null : ticketRepository.findById(commentDTO.getTicket())
                .orElseThrow(() -> new NotFoundException("ticket not found"));
        comment.setTicket(ticket);
        final User user = commentDTO.getUser() == null ? null : userRepository.findById(commentDTO.getUser())
                .orElseThrow(() -> new NotFoundException("user not found"));
        comment.setUser(user);
        return comment;
    }

    @EventListener(BeforeDeleteTicket.class)
    public void on(final BeforeDeleteTicket event) {
        final ReferencedException referencedException = new ReferencedException();
        final Comment ticketComment = commentRepository.findFirstByTicketId(event.getId());
        if (ticketComment != null) {
            referencedException.setKey("ticket.comment.ticket.referenced");
            referencedException.addParam(ticketComment.getId());
            throw referencedException;
        }
    }

    @EventListener(BeforeDeleteUser.class)
    public void on(final BeforeDeleteUser event) {
        final ReferencedException referencedException = new ReferencedException();
        final Comment userComment = commentRepository.findFirstByUserId(event.getId());
        if (userComment != null) {
            referencedException.setKey("user.comment.user.referenced");
            referencedException.addParam(userComment.getId());
            throw referencedException;
        }
    }

}

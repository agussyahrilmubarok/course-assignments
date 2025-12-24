package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.TicketComment;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.model.TicketCommentDTO;
import com.example.witrack.backend.repos.TicketCommentRepository;
import com.example.witrack.backend.repos.TicketRepository;
import com.example.witrack.backend.security.user.CurrentUserDetails;
import com.example.witrack.backend.service.TicketCommentService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Service
@Slf4j
@RequiredArgsConstructor
public class TicketCommentServiceImpl implements TicketCommentService {

    private final TicketCommentRepository ticketCommentRepository;
    private final TicketRepository ticketRepository;
    private final CurrentUserDetails currentUserDetails;

    @Override
    public TicketCommentDTO.TicketCommentResponse create(String code, TicketCommentDTO.TicketCommentRequest request) {
        Ticket ticket = ticketRepository.findByCode(code)
                .orElseThrow(() -> new NotFoundException("Ticket not found"));

        User currentUser = currentUserDetails.getUser();

        TicketComment ticketComment = new TicketComment();
        ticketComment.setContent(request.getContent());
        ticketComment.setUser(currentUser);
        ticketComment.setTicket(ticket);
        ticketComment = ticketCommentRepository.save(ticketComment);

        return TicketCommentDTO.TicketCommentResponse.fromTicketComment(ticketComment);
    }
}

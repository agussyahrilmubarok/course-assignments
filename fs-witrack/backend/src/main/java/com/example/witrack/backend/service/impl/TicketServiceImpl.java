package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.TicketReply;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.exception.UnauthorizedException;
import com.example.witrack.backend.model.*;
import com.example.witrack.backend.repository.TicketReplyRepository;
import com.example.witrack.backend.repository.TicketRepository;
import com.example.witrack.backend.repository.UserRepository;
import com.example.witrack.backend.security.CurrentUserDetails;
import com.example.witrack.backend.service.TicketService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.Clock;
import java.time.Instant;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
@Slf4j
@RequiredArgsConstructor
public class TicketServiceImpl implements TicketService {

    private final TicketRepository ticketRepository;
    private final TicketReplyRepository ticketReplyRepository;
    private final CurrentUserDetails currentUserDetails;
    private final UserRepository userRepository;

    @Override
    public List<TicketResponse> getTickets(String search, String status, String priority) {
        List<Ticket> tickets;

        if (search != null && !search.isEmpty()) {
            tickets = ticketRepository
                    .findByCodeIgnoreCaseContainingOrTitleIgnoreCaseContainingOrDescriptionIgnoreCaseContaining(
                            search, search, search
                    );
        } else if (status != null && priority != null) {
            tickets = ticketRepository.findByStatusAndPriority(status.toUpperCase(), priority.toUpperCase());
        } else if (status != null) {
            tickets = ticketRepository.findByStatus(status.toUpperCase());
        } else if (priority != null) {
            tickets = ticketRepository.findByPriority(priority.toUpperCase());
        } else {
            tickets = ticketRepository.findAll();
        }

        return tickets.stream()
                .map(TicketResponse::fromTicket)
                .collect(Collectors.toList());
    }

    @Override
    public TicketDetailResponse getTicketByCode(String code) {
        Ticket ticket = ticketRepository.findByCode(code)
                .orElseThrow(() -> new NotFoundException("Ticket code is not found: " + code));

        return TicketDetailResponse.fromTicket(ticket);
    }

    @Override
    public TicketResponse createTicket(TicketStoreRequest request) {
        User user = getUserById(currentUserDetails.getId());

        Ticket ticket = new Ticket();
        ticket.setTitle(request.getTitle());
        ticket.setDescription(request.getDescription());
        ticket.setStatus(Ticket.Status.valueOf(request.getStatus()));
        ticket.setPriority(Ticket.Priority.valueOf(request.getPriority()));
        ticket.setId(UUID.randomUUID().toString());
        ticket.setCode(generateCode());
        ticket.setCompleteAt(null);
        ticket.setUser(user);

        return TicketResponse.fromTicket(ticketRepository.save(ticket));
    }

    @Override
    public TicketResponse updateTicket(String code, TicketStoreRequest request) {
        Ticket ticket = ticketRepository.findByCode(code)
                .orElseThrow(() -> new NotFoundException("Ticket code is not found: " + code));
        User user = getUserById(currentUserDetails.getId());

        if (!ticket.getUser().getId().equals(user.getId()) && !user.hasRole(User.Role.ROLE_ADMIN)) {
            throw new UnauthorizedException("You are not allowed to update this ticket");
        }

        if (request.getTitle() != null) {
            ticket.setTitle(request.getTitle());
        }
        if (request.getDescription() != null) {
            ticket.setDescription(request.getDescription());
        }
        if (request.getStatus() != null) {
            ticket.setStatus(Ticket.Status.valueOf(request.getStatus()));
        }
        if (request.getPriority() != null) {
            ticket.setPriority(Ticket.Priority.valueOf(request.getPriority()));
        }

        return TicketResponse.fromTicket(ticketRepository.save(ticket));
    }

    @Override
    public void deleteTicketByCode(String code) {
        Ticket ticket = ticketRepository.findByCode(code)
                .orElseThrow(() -> new NotFoundException("Ticket code is not found: " + code));
        User user = getUserById(currentUserDetails.getId());

        if (!ticket.getUser().getId().equals(user.getId()) && !user.hasRole(User.Role.ROLE_ADMIN)) {
            throw new UnauthorizedException("You are not allowed to update this ticket");
        }

        ticketRepository.delete(ticket);
    }

    @Override
    public TicketReplyResponse createTicketReply(String code, TicketReplyStoreRequest request) {
        Ticket ticket = ticketRepository.findByCode(code)
                .orElseThrow(() -> new NotFoundException("Ticket code is not found: " + code));
        User user = getUserById(currentUserDetails.getId());

        TicketReply ticketReply = new TicketReply();
        ticketReply.setId(UUID.randomUUID().toString());
        ticketReply.setContent(request.getContent());
        ticketReply.setTicket(ticket);
        ticketReply.setUser(user);

        return TicketReplyResponse.fromTicketReply(ticketReplyRepository.save(ticketReply));
    }

    private User getUserById(String id) {
        return userRepository.findById(id)
                .orElseThrow(() -> new NotFoundException("User id is not found: " + id));
    }

    private String generateCode() {
        String hex = UUID.randomUUID().toString().replace("-", "");
        String first4 = hex.substring(0, 4);
        String last4 = hex.substring(hex.length() - 4);

        Clock clock = Clock.systemUTC();
        long unixTime = Instant.now(clock).getEpochSecond();

        return "TIC" + first4.toUpperCase() + unixTime + last4.toUpperCase();
    }
}

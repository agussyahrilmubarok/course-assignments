package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.exception.UnauthorizedException;
import com.example.witrack.backend.model.TicketDTO;
import com.example.witrack.backend.repos.TicketRepository;
import com.example.witrack.backend.security.user.CurrentUserDetails;
import com.example.witrack.backend.service.TicketService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Clock;
import java.time.Instant;
import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class TicketServiceImpl implements TicketService {

    private final TicketRepository ticketRepository;
    private final CurrentUserDetails currentUserDetails;

    @Override
    @Transactional
    public TicketDTO.TicketResponse create(TicketDTO.TicketRequest request) {
        User currentUser = currentUserDetails.getUser();

        Ticket ticket = new Ticket();
        ticket.setCode(generateTicketCode());
        ticket.setTitle(request.getTitle());
        ticket.setDescription(request.getDescription());
        ticket.setStatus(Ticket.Status.valueOf(request.getStatus()));
        ticket.setPriority(Ticket.Priority.valueOf(request.getPriority()));
        ticket.setCompleteAt(null);
        ticket.setUser(currentUser);
        ticket = ticketRepository.save(ticket);

        log.info("Ticket created successfully: code={}, userId={}, title={}", ticket.getCode(), currentUser.getId(), ticket.getTitle());
        return TicketDTO.TicketResponse.fromTicket(ticket);
    }

    @Override
    @Transactional
    public TicketDTO.TicketResponse update(UUID id, TicketDTO.TicketRequest request) {
        Ticket ticket = ticketRepository.findById(id)
                .orElseThrow(() -> new NotFoundException("Ticket not found"));
        validateTicketAccess(ticket);

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

        return TicketDTO.TicketResponse.fromTicket(ticketRepository.save(ticket));
    }

    @Override
    public void delete(UUID id) {
        Ticket ticket = ticketRepository.findById(id)
                .orElseThrow(() -> new NotFoundException("Ticket not found"));
        validateTicketAccess(ticket);
        ticketRepository.delete(ticket);
    }

    private String generateTicketCode() {
        String hex = UUID.randomUUID().toString().replace("-", "");
        String first4 = hex.substring(0, 4);
        String last4 = hex.substring(hex.length() - 4);
        Clock clock = Clock.systemUTC();
        long unixTime = Instant.now(clock).getEpochSecond();
        return "TIC" + first4.toUpperCase() + unixTime + last4.toUpperCase();
    }

    private void validateTicketAccess(Ticket ticket) {
        User currentUser = currentUserDetails.getUser();
        boolean isOwner = ticket.getUser().getId().equals(currentUser.getId());
        boolean isAdmin = currentUser.getRoles().contains(User.Role.ROLE_ADMIN);
        if (!isOwner && !isAdmin) {
            throw new UnauthorizedException("You do not have permission to access this ticket");
        }
    }
}

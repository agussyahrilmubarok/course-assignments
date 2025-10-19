package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.TicketReply;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.exception.UnauthorizedException;
import com.example.witrack.backend.model.*;
import com.example.witrack.backend.repository.TicketCustomRepository;
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
import java.time.OffsetDateTime;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
@Slf4j
@RequiredArgsConstructor
public class TicketServiceImpl implements TicketService {

    private final TicketRepository ticketRepository;
    private final TicketReplyRepository ticketReplyRepository;
    private final TicketCustomRepository ticketCustomRepository;
    private final CurrentUserDetails currentUserDetails;
    private final UserRepository userRepository;

    @Override
    public List<TicketResponse> getTickets(String keyword, String status, String priority, String date) {
        Ticket.Status statusEnum = null;
        if (status != null && !status.isBlank()) {
            statusEnum = Ticket.Status.valueOf(status.toUpperCase());
        }

        Ticket.Priority priorityEnum = null;
        if (priority != null && !priority.isBlank()) {
            priorityEnum = Ticket.Priority.valueOf(priority.toUpperCase());
        }

        OffsetDateTime now = OffsetDateTime.now();
        OffsetDateTime startAt = null;
        OffsetDateTime endAt = null;

        if ("TODAY".equalsIgnoreCase(date)) {
            startAt = now.toLocalDate().atStartOfDay().atOffset(now.getOffset());
            endAt = startAt.plusDays(1).minusNanos(1);
        } else if ("MONTH".equalsIgnoreCase(date)) {
            startAt = now.withDayOfMonth(1).toLocalDate().atStartOfDay().atOffset(now.getOffset());
            endAt = startAt.plusMonths(1).minusNanos(1);
        } else if ("YEAR".equalsIgnoreCase(date)) {
            startAt = now.withDayOfYear(1).toLocalDate().atStartOfDay().atOffset(now.getOffset());
            endAt = startAt.plusYears(1).minusNanos(1);
        }

        List<Ticket> tickets = ticketCustomRepository.searchTickets(
                keyword != null ? keyword : "",
                statusEnum,
                priorityEnum,
                startAt,
                endAt
        );

        return tickets.stream()
                .map(TicketResponse::fromTicket)
                .collect(Collectors.toList());
    }

    @Override
    public List<TicketResponse> getMyTickets(String keyword, String status, String priority, String date) {
        User user = getUserById(currentUserDetails.getId());

        Ticket.Status statusEnum = null;
        if (status != null && !status.isBlank()) {
            statusEnum = Ticket.Status.valueOf(status.toUpperCase());
        }

        Ticket.Priority priorityEnum = null;
        if (priority != null && !priority.isBlank()) {
            priorityEnum = Ticket.Priority.valueOf(priority.toUpperCase());
        }

        OffsetDateTime now = OffsetDateTime.now();
        OffsetDateTime startAt = null;
        OffsetDateTime endAt = null;

        if ("TODAY".equalsIgnoreCase(date)) {
            startAt = now.toLocalDate().atStartOfDay().atOffset(now.getOffset());
            endAt = startAt.plusDays(1).minusNanos(1);
        } else if ("MONTH".equalsIgnoreCase(date)) {
            startAt = now.withDayOfMonth(1).toLocalDate().atStartOfDay().atOffset(now.getOffset());
            endAt = startAt.plusMonths(1).minusNanos(1);
        } else if ("YEAR".equalsIgnoreCase(date)) {
            startAt = now.withDayOfYear(1).toLocalDate().atStartOfDay().atOffset(now.getOffset());
            endAt = startAt.plusYears(1).minusNanos(1);
        }

        List<Ticket> tickets = ticketCustomRepository.searchTickets(
                keyword != null ? keyword : "",
                statusEnum,
                priorityEnum,
                startAt,
                endAt,
                user.getId()
        );

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

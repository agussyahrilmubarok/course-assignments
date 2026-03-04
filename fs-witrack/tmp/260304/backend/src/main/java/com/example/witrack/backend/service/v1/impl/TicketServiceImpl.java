package com.example.witrack.backend.service.v1.impl;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.exception.UnauthorizedException;
import com.example.witrack.backend.model.TicketDTO;
import com.example.witrack.backend.repos.TicketRepository;
import com.example.witrack.backend.security.user.CurrentUserDetails;
import com.example.witrack.backend.service.v1.TicketService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.*;
import java.util.List;
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
    @Transactional
    public void delete(UUID id) {
        Ticket ticket = ticketRepository.findById(id)
                .orElseThrow(() -> new NotFoundException("Ticket not found"));
        validateTicketAccess(ticket);
        ticketRepository.delete(ticket);
    }

    @Override
    @Transactional
    public TicketDTO.TicketResponse findById(UUID id) {
        Ticket ticket = ticketRepository.findById(id)
                .orElseThrow(() -> new NotFoundException("Ticket not found"));
        return TicketDTO.TicketResponse.fromTicket(ticket);
    }

    @Override
    @Transactional
    public TicketDTO.TicketResponse findByCode(String code) {
        Ticket ticket = ticketRepository.findByCode(code)
                .orElseThrow(() -> new NotFoundException("Ticket not found"));
        return TicketDTO.TicketResponse.fromTicket(ticket);
    }

    @Override
    @Transactional(readOnly = true)
    public List<TicketDTO.TicketResponse> searchTicket(String search, String status, String priority, String date) {
        User currentUser = currentUserDetails.getUser();
        boolean isAdmin = currentUser.getRoles().contains(User.Role.ROLE_ADMIN);

        Specification<Ticket> spec = buildSearchSpecification(search, status, priority, date);

        if (!isAdmin) {
            spec = spec.and((root, query, cb) ->
                    cb.equal(root.get("user").get("id"), currentUser.getId()));
        }

        return ticketRepository.findAll(spec)
                .stream()
                .map(TicketDTO.TicketResponse::fromTicket)
                .toList();
    }

    @Override
    public List<TicketDTO.TicketResponse> searchMyTicket(String search, String status, String priority, String date) {
        User currentUser = currentUserDetails.getUser();

        Specification<Ticket> spec = buildSearchSpecification(search, status, priority, date)
                .and((root, query, cb) -> cb.equal(root.get("user").get("id"), currentUser.getId()));

        return ticketRepository.findAll(spec)
                .stream()
                .map(TicketDTO.TicketResponse::fromTicket)
                .toList();
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

    private Specification<Ticket> buildSearchSpecification(String search, String status,
                                                           String priority, String date) {
        return (root, query, cb) -> {
            var predicates = cb.conjunction();

            if (search != null && !search.isBlank()) {
                String like = "%" + search.toLowerCase() + "%";
                predicates = cb.and(predicates,
                        cb.or(
                                cb.like(cb.lower(root.get("title")), like),
                                cb.like(cb.lower(root.get("description")), like),
                                cb.like(cb.lower(root.get("code")), like)
                        )
                );
            }

            if (status != null) {
                predicates = cb.and(predicates, cb.equal(root.get("status"), Ticket.Status.valueOf(status)));
            }

            if (priority != null) {
                predicates = cb.and(predicates, cb.equal(root.get("priority"), Ticket.Priority.valueOf(priority)));
            }

            if (date != null) {
                LocalDate localDate = LocalDate.parse(date);
                OffsetDateTime start = localDate.atStartOfDay().atOffset(ZoneOffset.UTC);
                OffsetDateTime end = start.plusDays(1);

                predicates = cb.and(predicates, cb.between(root.get("createdAt"), start, end));
            }

            return predicates;
        };
    }
}

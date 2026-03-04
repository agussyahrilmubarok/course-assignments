package com.example.witrack.backend.service;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.events.BeforeDeleteTicket;
import com.example.witrack.backend.events.BeforeDeleteUser;
import com.example.witrack.backend.model.TicketDTO;
import com.example.witrack.backend.repos.TicketRepository;
import com.example.witrack.backend.repos.UserRepository;
import com.example.witrack.backend.util.NotFoundException;
import com.example.witrack.backend.util.ReferencedException;
import java.util.List;
import java.util.UUID;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.context.event.EventListener;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;


@Service
public class TicketService {

    private final TicketRepository ticketRepository;
    private final UserRepository userRepository;
    private final ApplicationEventPublisher publisher;

    public TicketService(final TicketRepository ticketRepository,
            final UserRepository userRepository, final ApplicationEventPublisher publisher) {
        this.ticketRepository = ticketRepository;
        this.userRepository = userRepository;
        this.publisher = publisher;
    }

    public List<TicketDTO> findAll() {
        final List<Ticket> tickets = ticketRepository.findAll(Sort.by("id"));
        return tickets.stream()
                .map(ticket -> mapToDTO(ticket, new TicketDTO()))
                .toList();
    }

    public TicketDTO get(final UUID id) {
        return ticketRepository.findById(id)
                .map(ticket -> mapToDTO(ticket, new TicketDTO()))
                .orElseThrow(NotFoundException::new);
    }

    public UUID create(final TicketDTO ticketDTO) {
        final Ticket ticket = new Ticket();
        mapToEntity(ticketDTO, ticket);
        return ticketRepository.save(ticket).getId();
    }

    public void update(final UUID id, final TicketDTO ticketDTO) {
        final Ticket ticket = ticketRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        mapToEntity(ticketDTO, ticket);
        ticketRepository.save(ticket);
    }

    public void delete(final UUID id) {
        final Ticket ticket = ticketRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        publisher.publishEvent(new BeforeDeleteTicket(id));
        ticketRepository.delete(ticket);
    }

    private TicketDTO mapToDTO(final Ticket ticket, final TicketDTO ticketDTO) {
        ticketDTO.setId(ticket.getId());
        ticketDTO.setCode(ticket.getCode());
        ticketDTO.setTitle(ticket.getTitle());
        ticketDTO.setDescription(ticket.getDescription());
        ticketDTO.setStatus(ticket.getStatus());
        ticketDTO.setPriority(ticket.getPriority());
        ticketDTO.setCompleteAt(ticket.getCompleteAt());
        ticketDTO.setUser(ticket.getUser() == null ? null : ticket.getUser().getId());
        return ticketDTO;
    }

    private Ticket mapToEntity(final TicketDTO ticketDTO, final Ticket ticket) {
        ticket.setCode(ticketDTO.getCode());
        ticket.setTitle(ticketDTO.getTitle());
        ticket.setDescription(ticketDTO.getDescription());
        ticket.setStatus(ticketDTO.getStatus());
        ticket.setPriority(ticketDTO.getPriority());
        ticket.setCompleteAt(ticketDTO.getCompleteAt());
        final User user = ticketDTO.getUser() == null ? null : userRepository.findById(ticketDTO.getUser())
                .orElseThrow(() -> new NotFoundException("user not found"));
        ticket.setUser(user);
        return ticket;
    }

    public boolean codeExists(final String code) {
        return ticketRepository.existsByCodeIgnoreCase(code);
    }

    @EventListener(BeforeDeleteUser.class)
    public void on(final BeforeDeleteUser event) {
        final ReferencedException referencedException = new ReferencedException();
        final Ticket userTicket = ticketRepository.findFirstByUserId(event.getId());
        if (userTicket != null) {
            referencedException.setKey("user.ticket.user.referenced");
            referencedException.addParam(userTicket.getId());
            throw referencedException;
        }
    }

}

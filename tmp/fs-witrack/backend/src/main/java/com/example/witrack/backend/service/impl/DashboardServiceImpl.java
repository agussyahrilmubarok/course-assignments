package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.model.DashboardResponse;
import com.example.witrack.backend.repository.TicketRepository;
import com.example.witrack.backend.service.DashboardService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.time.OffsetDateTime;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Service
@Slf4j
@RequiredArgsConstructor
public class DashboardServiceImpl implements DashboardService {

    private final TicketRepository ticketRepository;

    @Override
    public DashboardResponse getStatistics() {
        OffsetDateTime now = OffsetDateTime.now();
        OffsetDateTime startOfMonth = now.withDayOfMonth(1)
                .withHour(0).withMinute(0).withSecond(0).withNano(0);
        OffsetDateTime endOfMonth = startOfMonth.plusMonths(1).minusNanos(1);

        List<Ticket> tickets = ticketRepository.findByCreatedAtBetween(startOfMonth, endOfMonth);

        long totalTickets = tickets.size();
        long resolvedTickets = tickets.stream()
                .filter(t -> t.getStatus() == Ticket.Status.RESOLVED)
                .count();
        long activeTickets = tickets.stream()
                .filter(t -> t.getStatus() != Ticket.Status.RESOLVED)
                .count();

        double avgResolutionTime = tickets.stream()
                .filter(t -> t.getStatus() == Ticket.Status.RESOLVED && t.getCompleteAt() != null)
                .mapToLong(t -> Duration.between(t.getCreatedAt(), t.getCompleteAt()).toHours())
                .average()
                .orElse(0);

        Map<String, Long> statusDistribution = tickets.stream()
                .collect(Collectors.groupingBy(
                        t -> t.getStatus().name().toLowerCase(),
                        Collectors.counting()));

        return DashboardResponse.builder()
                .totalTickets(totalTickets)
                .activeTickets(activeTickets)
                .resolvedTickets(resolvedTickets)
                .avgResolutionTime(Math.round(avgResolutionTime * 10.0) / 10.0) // 1 decimal
                .statusDistribution(statusDistribution)
                .build();
    }
}

package com.example.witrack.backend.model;

import lombok.Builder;
import lombok.Data;

import java.util.Map;

@Data
@Builder
public class DashboardResponse {

    private long totalTickets;

    private long activeTickets;

    private long resolvedTickets;

    private double avgResolutionTime;

    private Map<String, Long> statusDistribution;
}

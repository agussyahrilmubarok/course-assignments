package com.example.witrack.backend.model;

import java.util.Map;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class DashboardResponse {

    private long totalTickets;

    private long activeTickets;
    
    private long resolvedTickets;
    
    private double avgResolutionTime;
    
    private Map<String, Long> statusDistribution;
}

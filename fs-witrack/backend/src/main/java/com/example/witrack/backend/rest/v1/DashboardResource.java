package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.DashboardResponse;
import com.example.witrack.backend.service.DashboardService;

import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController("DashboardResourceV1")
@RequestMapping(value = "/api/v1/dashboards", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class DashboardResource {

    private final DashboardService dashboardService;

    @GetMapping("/statistics")
    public ResponseEntity<DashboardResponse> getStatistics() {
        DashboardResponse response = dashboardService.getStatistics();
        return ResponseEntity.ok(response);
    }
}

package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.*;
import com.example.witrack.backend.service.TicketService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController("TicketResourceV1")
@RequestMapping(value = "/api/v1/tickets", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class TicketResource {

    private final TicketService ticketService;

    @GetMapping
    public ResponseEntity<List<TicketResponse>> getTickets(@RequestParam(required = false) String search,
                                                           @RequestParam(required = false) String status,
                                                           @RequestParam(required = false) String priority,
                                                           @RequestParam(required = false) String date) {
        List<TicketResponse> responses = ticketService.getTickets(search, status, priority, date);
        return ResponseEntity.ok(responses);
    }

    @GetMapping("/me")
    public ResponseEntity<List<TicketResponse>> getMyTickets(@RequestParam(required = false) String search,
                                                             @RequestParam(required = false) String status,
                                                             @RequestParam(required = false) String priority,
                                                             @RequestParam(required = false) String date) {
        List<TicketResponse> responses = ticketService.getMyTickets(search, status, priority, date);
        return ResponseEntity.ok(responses);
    }

    @GetMapping("/{code}")
    public ResponseEntity<TicketDetailResponse> getTicketByCode(@PathVariable("code") String code) {
        TicketDetailResponse response = ticketService.getTicketByCode(code);
        return ResponseEntity.ok(response);
    }

    @PostMapping
    public ResponseEntity<TicketResponse> createTicket(@Valid @RequestBody TicketStoreRequest request) {
        TicketResponse response = ticketService.createTicket(request);
        return ResponseEntity.ok(response);
    }

    @PutMapping("/{code}")
    public ResponseEntity<TicketResponse> updateTicketByCode(@PathVariable("code") String code,
                                                             @Valid @RequestBody TicketStoreRequest request) {
        TicketResponse response = ticketService.updateTicket(code, request);
        return ResponseEntity.ok(response);
    }

    @DeleteMapping("/{code}")
    public ResponseEntity<Void> deleteTicketByCode(@PathVariable("code") String code) {
        ticketService.deleteTicketByCode(code);
        return ResponseEntity.noContent().build();
    }

    @PostMapping("/{code}/reply")
    public ResponseEntity<TicketReplyResponse> createTicketReply(@PathVariable("code") String code,
                                                                 @Valid @RequestBody TicketReplyStoreRequest request) {
        TicketReplyResponse response = ticketService.createTicketReply(code, request);
        return ResponseEntity.ok(response);
    }
}

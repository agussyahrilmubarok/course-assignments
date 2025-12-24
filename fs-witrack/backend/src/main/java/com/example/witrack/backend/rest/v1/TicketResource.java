package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.TicketCommentDTO;
import com.example.witrack.backend.model.TicketDTO;
import com.example.witrack.backend.service.v1.TicketCommentService;
import com.example.witrack.backend.service.v1.TicketService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController("TicketResourceV1")
@RequestMapping(value = "/api/v1/tickets", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class TicketResource {

    private final TicketService ticketService;
    private final TicketCommentService ticketCommentService;

    @PostMapping
    public ResponseEntity<TicketDTO.TicketResponse> create(@Valid @RequestBody TicketDTO.TicketRequest request) {
        return ResponseEntity.ok(ticketService.create(request));
    }

    @PutMapping("/{id}")
    public ResponseEntity<TicketDTO.TicketResponse> update(@PathVariable("id") UUID id,
                                                           @Valid @RequestBody TicketDTO.TicketRequest request) {
        return ResponseEntity.ok(ticketService.update(id, request));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable("id") UUID id) {
        ticketService.delete(id);
        return ResponseEntity.noContent().build();
    }

    @GetMapping("/{id}")
    public ResponseEntity<TicketDTO.TicketResponse> getById(@PathVariable("id") UUID id) {
        return ResponseEntity.ok(ticketService.findById(id));
    }

    @GetMapping("/code/{code}")
    public ResponseEntity<TicketDTO.TicketResponse> getByCode(@PathVariable("code") String code) {
        return ResponseEntity.ok(ticketService.findByCode(code));
    }

    @GetMapping
    public ResponseEntity<List<TicketDTO.TicketResponse>> searchTicket(@RequestParam(required = false) String search,
                                                                       @RequestParam(required = false) String status,
                                                                       @RequestParam(required = false) String priority,
                                                                       @RequestParam(required = false) String date) {
        return ResponseEntity.ok(ticketService.searchTicket(search, status, priority, date));
    }

    @GetMapping("/me")
    public ResponseEntity<List<TicketDTO.TicketResponse>> searchMyTicket(@RequestParam(required = false) String search,
                                                                         @RequestParam(required = false) String status,
                                                                         @RequestParam(required = false) String priority,
                                                                         @RequestParam(required = false) String date) {
        return ResponseEntity.ok(ticketService.searchMyTicket(search, status, priority, date));
    }

    @PostMapping("/{code}/comments")
    public ResponseEntity<TicketCommentDTO.TicketCommentResponse> createComment(@PathVariable("code") String code,
                                                                                @Valid @RequestBody TicketCommentDTO.TicketCommentRequest request) {
        return ResponseEntity.ok(ticketCommentService.create(code, request));
    }
}

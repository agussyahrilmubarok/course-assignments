package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.TicketDTO;
import com.example.witrack.backend.service.TicketService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController("TicketResourceV1")
@RequestMapping(value = "/api/v1/tickets", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class TicketResource {

    private final TicketService ticketService;

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
}

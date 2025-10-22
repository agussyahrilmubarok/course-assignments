package com.example.payment.rest;

import com.example.payment.model.TransactionDTO;
import com.example.payment.service.TransactionService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController("TransactionResource")
@RequestMapping(value = "/api/v1/transactions", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class TransactionResource {

    private final TransactionService transactionService;

    @GetMapping("/{id}")
    public ResponseEntity<TransactionDTO.Response> getTransactionById(@PathVariable String id) {
        TransactionDTO.Response response = transactionService.findById(id);
        return ResponseEntity.ok(response);
    }

    @PostMapping(consumes = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<TransactionDTO.Response> createTransaction(@RequestBody @Valid TransactionDTO.CreateTransactionRequest request) {
        TransactionDTO.Response created = transactionService.create(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(created);
    }
}

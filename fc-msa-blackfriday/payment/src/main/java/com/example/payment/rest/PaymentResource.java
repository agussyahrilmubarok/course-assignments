package com.example.payment.rest;

import com.example.payment.model.PaymentDTO;
import com.example.payment.service.PaymentService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping(value = "/api/v1/payments", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class PaymentResource {

    private final PaymentService paymentService;

    @GetMapping("/{id}")
    public ResponseEntity<PaymentDTO.Response> getTransactionById(@PathVariable String id) {
        PaymentDTO.Response response = paymentService.findById(id);
        return ResponseEntity.ok(response);
    }

    @PostMapping(consumes = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<PaymentDTO.Response> createTransaction(@RequestBody @Valid PaymentDTO.CreateTransactionRequest request) {
        PaymentDTO.Response created = paymentService.create(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(created);
    }
}

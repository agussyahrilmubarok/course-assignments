package com.example.order.service;

import com.example.order.model.PaymentDTO;
import jakarta.validation.Valid;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

@FeignClient(name = "PAYMENT-SERVICE")
public interface PaymentClient {

    @PostMapping("/api/v1/payments")
    PaymentDTO.Response createPayment(@RequestBody @Valid PaymentDTO.CreateTransactionRequest request);
}

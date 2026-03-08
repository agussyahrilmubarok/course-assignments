package com.example.payment.service;

import com.example.payment.model.PaymentDTO;

public interface PaymentService {

    PaymentDTO.Response findById(String id);

    PaymentDTO.Response create(PaymentDTO.CreateTransactionRequest param);
}

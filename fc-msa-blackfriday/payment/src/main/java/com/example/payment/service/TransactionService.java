package com.example.payment.service;

import com.example.payment.model.TransactionDTO;

public interface TransactionService {

    TransactionDTO.Response findById(String id);

    TransactionDTO.Response create(TransactionDTO.CreateTransactionRequest param);
}

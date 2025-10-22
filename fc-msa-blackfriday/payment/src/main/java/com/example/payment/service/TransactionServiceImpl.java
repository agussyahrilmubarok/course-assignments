package com.example.payment.service;

import com.example.payment.domain.Transaction;
import com.example.payment.exception.MidtransPaymentException;
import com.example.payment.exception.TransactionNotFoundException;
import com.example.payment.model.TransactionDTO;
import com.example.payment.repos.TransactionRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class TransactionServiceImpl implements TransactionService {

    private final TransactionRepository transactionRepository;
    private final MidtransService midtransService;

    @Override
    public TransactionDTO.Response findById(String id) {
        Transaction transaction = transactionRepository.findById(id)
                .orElseThrow(() -> new TransactionNotFoundException("Transaction not found with ID: " + id));
        return TransactionDTO.Response.from(transaction);
    }

    @Override
    @Transactional
    public TransactionDTO.Response create(TransactionDTO.CreateTransactionRequest param) {
        String transactionId = UUID.randomUUID().toString();

        Transaction transaction = new Transaction();
        transaction.setId(transactionId);
        transaction.setOrderId(param.getOrderId());
        transaction.setAmount(param.getAmount());
        transaction.setStatus(Transaction.Status.PENDING);
        transactionRepository.save(transaction);

        try {
            String paymentUrl = midtransService.createPaymentRedirectUrl(transactionId, param.getAmount());
            transaction.setPaymentUrl(paymentUrl);
            transactionRepository.save(transaction);
            log.info("Transaction created successfully with ID: {}", transactionId);
            return TransactionDTO.Response.from(transaction);

        } catch (MidtransPaymentException ex) {
            log.error("Midtrans error for transaction ID {}: {}", transactionId, ex.getMessage());
            transaction.setStatus(Transaction.Status.FAILED);
            transactionRepository.save(transaction);
            throw ex;
        }
    }
}

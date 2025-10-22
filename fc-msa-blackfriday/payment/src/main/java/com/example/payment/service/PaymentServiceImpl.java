package com.example.payment.service;

import com.example.payment.domain.Payment;
import com.example.payment.exception.MidtransPaymentException;
import com.example.payment.exception.TransactionNotFoundException;
import com.example.payment.model.PaymentDTO;
import com.example.payment.repos.PaymentRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class PaymentServiceImpl implements PaymentService {

    private final PaymentRepository paymentRepository;
    private final MidtransService midtransService;

    @Override
    public PaymentDTO.Response findById(String id) {
        Payment payment = paymentRepository.findById(id)
                .orElseThrow(() -> new TransactionNotFoundException("Transaction not found with ID: " + id));
        return PaymentDTO.Response.from(payment);
    }

    @Override
    @Transactional
    public PaymentDTO.Response create(PaymentDTO.CreateTransactionRequest param) {
        String transactionId = UUID.randomUUID().toString();

        Payment payment = new Payment();
        payment.setId(transactionId);
        payment.setOrderId(param.getOrderId());
        payment.setAmount(param.getAmount());
        payment.setStatus(Payment.Status.PENDING);
        paymentRepository.save(payment);

        try {
            String paymentUrl = midtransService.createPaymentRedirectUrl(transactionId, param.getAmount());
            payment.setPaymentUrl(paymentUrl);
            paymentRepository.save(payment);
            log.info("Transaction created successfully with ID: {}", transactionId);
            return PaymentDTO.Response.from(payment);

        } catch (MidtransPaymentException ex) {
            log.error("Midtrans error for transaction ID {}: {}", transactionId, ex.getMessage());
            payment.setStatus(Payment.Status.FAILED);
            paymentRepository.save(payment);
            throw ex;
        }
    }
}

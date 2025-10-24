package com.example.payment.service;

import java.math.BigDecimal;

public interface MidtransService {

    String createPaymentRedirectUrl(String transactionId, BigDecimal amount);
}

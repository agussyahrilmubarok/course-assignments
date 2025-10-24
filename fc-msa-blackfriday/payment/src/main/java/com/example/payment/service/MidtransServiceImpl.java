package com.example.payment.service;

import com.example.payment.exception.MidtransPaymentException;
import com.midtrans.Config;
import com.midtrans.ConfigFactory;
import com.midtrans.httpclient.error.MidtransError;
import com.midtrans.service.MidtransSnapApi;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.util.HashMap;
import java.util.Map;

@Service
@Slf4j
@RequiredArgsConstructor
public class MidtransServiceImpl implements MidtransService {

    @Value("${midtrans.server-key}")
    private String midtransServerKey;

    @Value("${midtrans.client-key}")
    private String midtransClientKey;

    @Value("${midtrans.is-production:false}")
    private boolean midtransIsProduction;

    @Override
    public String createPaymentRedirectUrl(String transactionId, BigDecimal amount) {
        MidtransSnapApi snapApi = new ConfigFactory(
                new Config(midtransServerKey, midtransClientKey, midtransIsProduction)
        ).getSnapApi();

        Map<String, Object> params = new HashMap<>();
        Map<String, String> transactionDetails = new HashMap<>();
        transactionDetails.put("order_id", transactionId);
        transactionDetails.put("gross_amount", amount.toString());
        params.put("transaction_details", transactionDetails);

        try {
            return snapApi.createTransactionRedirectUrl(params);
        } catch (MidtransError e) {
            log.error("Failed to create Midtrans payment URL: {}", e.getMessage(), e);
            throw new MidtransPaymentException("Failed to create Midtrans payment URL", e);
        }
    }
}

package com.example.payment.exception;

public class MidtransPaymentException extends RuntimeException {

    public MidtransPaymentException(String message) {
        super(message);
    }

    public MidtransPaymentException(String message, Throwable cause) {
        super(message, cause);
    }
}

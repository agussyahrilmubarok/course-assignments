package com.example.coupon.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.NOT_FOUND)
public class CouponNotFoundException extends RuntimeException {

    public CouponNotFoundException() {
        super();
    }

    public CouponNotFoundException(final String message) {
        super(message);
    }
}

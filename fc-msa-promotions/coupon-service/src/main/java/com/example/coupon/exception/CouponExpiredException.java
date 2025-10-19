package com.example.coupon.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.GONE)
public class CouponExpiredException extends RuntimeException {

    public CouponExpiredException() {
        super();
    }

    public CouponExpiredException(final String message) {
        super(message);
    }
}

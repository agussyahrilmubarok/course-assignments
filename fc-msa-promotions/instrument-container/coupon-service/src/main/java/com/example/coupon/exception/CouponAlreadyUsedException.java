package com.example.coupon.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.CONFLICT)
public class CouponAlreadyUsedException extends RuntimeException {

    public CouponAlreadyUsedException() {
        super();
    }

    public CouponAlreadyUsedException(final String message) {
        super(message);
    }
}

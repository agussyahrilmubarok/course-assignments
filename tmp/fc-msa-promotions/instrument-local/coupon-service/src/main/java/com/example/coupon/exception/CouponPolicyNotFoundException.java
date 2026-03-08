package com.example.coupon.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.NOT_FOUND)
public class CouponPolicyNotFoundException extends RuntimeException {

    public CouponPolicyNotFoundException() {
        super();
    }

    public CouponPolicyNotFoundException(final String message) {
        super(message);
    }
}

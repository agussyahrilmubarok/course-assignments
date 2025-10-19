package com.example.coupon.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.UNPROCESSABLE_ENTITY)
public class CouponNotUsedException extends RuntimeException {

    public CouponNotUsedException(String message) {
        super(message);
    }

    public CouponNotUsedException(String message, Throwable cause) {
        super(message, cause);
    }
}

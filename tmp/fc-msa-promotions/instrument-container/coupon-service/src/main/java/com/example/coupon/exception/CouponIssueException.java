package com.example.coupon.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.BAD_REQUEST)
public class CouponIssueException extends RuntimeException {

    public CouponIssueException() {
        super();
    }

    public CouponIssueException(final String message) {
        super(message);
    }
}

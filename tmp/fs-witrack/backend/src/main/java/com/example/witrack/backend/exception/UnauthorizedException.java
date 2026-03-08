package com.example.witrack.backend.exception;


import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.FORBIDDEN)
public class UnauthorizedException extends RuntimeException {

    public UnauthorizedException() {
        super();
    }

    public UnauthorizedException(final String message) {
        super(message);
    }
}
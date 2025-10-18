package com.example.gateway.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.server.ResponseStatusException;

public class InvalidTokenException extends ResponseStatusException {

    public InvalidTokenException() {
        super(HttpStatus.UNAUTHORIZED, "Invalid or expired token");
    }

    public InvalidTokenException(String reason) {
        super(HttpStatus.UNAUTHORIZED, reason);
    }
}

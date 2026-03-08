package com.example.user.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.BAD_REQUEST)
public class UnauthorizedAccessException extends RuntimeException {

    public UnauthorizedAccessException() {
        super();
    }

    public UnauthorizedAccessException(final String message) {
        super(message);
    }
}

package com.example.user.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.BAD_REQUEST)
public class DuplicateUserException extends RuntimeException {

    public DuplicateUserException() {
        super();
    }

    public DuplicateUserException(final String message) {
        super(message);
    }
}

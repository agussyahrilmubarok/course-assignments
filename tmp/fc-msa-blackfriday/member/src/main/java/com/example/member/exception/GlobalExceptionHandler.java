package com.example.member.exception;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.time.LocalDateTime;
import java.util.Map;

@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(UserAlreadyExistsException.class)
    public ResponseEntity<?> handleUserAlreadyExists(UserAlreadyExistsException ex) {
        return ResponseEntity.status(HttpStatus.CONFLICT).body(
                Map.of("error", ex.getMessage(), "status", 409, "timestamp", LocalDateTime.now())
        );
    }

    @ExceptionHandler(UserNotFoundException.class)
    public ResponseEntity<?> handleUserNotFound(UserNotFoundException ex) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(
                Map.of("error", ex.getMessage(), "status", 404, "timestamp", LocalDateTime.now())
        );
    }

    @ExceptionHandler(InvalidPasswordException.class)
    public ResponseEntity<?> handleInvalidPassword(InvalidPasswordException ex) {
        return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(
                Map.of("error", ex.getMessage(), "status", 401, "timestamp", LocalDateTime.now())
        );
    }
}

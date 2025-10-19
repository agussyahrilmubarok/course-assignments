package com.example.order.util;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
public class OrderEventPublishException extends RuntimeException {

    public OrderEventPublishException() {
        super();
    }

    public OrderEventPublishException(final String message) {
        super(message);
    }

    public OrderEventPublishException(final String message, final Exception exception) {
        super(message, exception);
    }
}

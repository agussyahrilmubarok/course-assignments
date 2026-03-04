package com.example.witrack.backend.exception;

import com.example.witrack.backend.util.ReferencedWarning;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.CONFLICT)
public class ReferencedException extends RuntimeException {

    public ReferencedException() {
        super();
    }

    public ReferencedException(final ReferencedWarning referencedWarning) {
        super(referencedWarning.toMessage());
    }
}

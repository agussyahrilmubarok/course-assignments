package com.example.witrack.backend.security;

import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;

@Component
@Slf4j
public class CurrentUserDetails {

    public String getId() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth == null || "anonymousUser".equals(auth.getPrincipal())) {
            log.warn("Authentication attempt failed or no authenticated user found.");
            throw new BadCredentialsException("User is not authenticated.");
        }

        Object principal = auth.getPrincipal();
        if (principal instanceof UserDetailsImpl userDetails) {
            return userDetails.getId();
        }
        throw new BadCredentialsException("User is not authenticated.");
    }
}

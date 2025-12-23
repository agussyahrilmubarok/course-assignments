package com.example.witrack.backend.security.user;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.repos.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;

import java.util.UUID;

@Component
@RequiredArgsConstructor
public class CurrentUserDetails {

    private final UserRepository userRepository;

    public UUID getId() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth == null || "anonymousUser".equals(auth.getPrincipal())) {
            throw new BadCredentialsException("User is not authenticated.");
        }

        Object principal = auth.getPrincipal();
        if (principal instanceof UserDetailsImpl userDetails) {
            return userDetails.getId();
        }

        throw new BadCredentialsException("User is not authenticated.");
    }

    public User getUser() {
        UUID userId = getId();
        return userRepository.findById(userId)
                .orElseThrow(() -> new BadCredentialsException("Authenticated user not found in database."));
    }
}

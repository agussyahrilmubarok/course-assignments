package com.finly.finly.security;

import com.finly.finly.domain.User;
import com.finly.finly.repos.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

@Service
@Slf4j
@RequiredArgsConstructor
public class UserDetailsServiceImpl implements UserDetailsService {

    private final UserRepository userRepository;

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        log.info("Loading user details for email: {}", username);

        User user = userRepository.findByEmail(username)
                .orElseThrow(() -> {
                    log.warn("No user found with email: {}", username);
                    return new UsernameNotFoundException("Invalid email or password");
                });

        log.info("User details successfully loaded for email: {}", username);
        return new UserDetailsImpl(user);
    }
}

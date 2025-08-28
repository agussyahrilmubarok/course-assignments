package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.model.UserResponse;
import com.example.witrack.backend.repository.UserRepository;
import com.example.witrack.backend.service.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Service
@Slf4j
@RequiredArgsConstructor
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;

    @Override
    public UserResponse getById(String id) {
        User user = userRepository.findById(id)
                .orElseThrow(() -> {
                    log.warn("User not found with id: {}", id);
                    return new NotFoundException("User id is not found: " + id);
                });

        log.info("User found with id: {} and email: {}", user.getId(), user.getEmail());
        return UserResponse.fromUser(user);
    }
}

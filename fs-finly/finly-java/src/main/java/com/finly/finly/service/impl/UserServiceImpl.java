package com.finly.finly.service.impl;


import com.finly.finly.domain.User;
import com.finly.finly.model.UserDTO;
import com.finly.finly.repos.UserRepository;
import com.finly.finly.service.UserService;
import com.finly.finly.util.NotFoundException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;

    @Override
    public List<UserDTO> findAll() {
        log.info("Retrieving all users from the database.");
        return userRepository.findAll()
                .stream()
                .map(UserDTO::fromUser)
                .toList();
    }

    @Override
    public UserDTO get(UUID id) {
        log.info("Retrieving user details for ID: {}", id);
        return userRepository.findById(id)
                .map(UserDTO::fromUser)
                .orElseThrow(() -> {
                    log.error("User with ID {} was not found.", id);
                    return new NotFoundException("User not found for ID: " + id);
                });
    }

    @Override
    public UUID create(UserDTO userDTO) {
        log.info("Creating a new user with email: {}", userDTO.getEmail());
        User user = new User();
        user.setEmail(userDTO.getEmail());
        user.setPassword(passwordEncoder.encode(userDTO.getPassword()));
        UUID savedId = userRepository.save(user).getId();
        log.info("User created successfully with ID: {}", savedId);
        return savedId;
    }

    @Override
    public void update(UUID id, UserDTO userDTO) {
        log.info("Updating user information for ID: {}", id);
        User user = userRepository.findById(id)
                .orElseThrow(() -> {
                    log.error("User with ID {} not found. Update aborted.", id);
                    return new NotFoundException("User not found for ID: " + id);
                });
        user.setEmail(userDTO.getEmail());
        user.setPassword(passwordEncoder.encode(userDTO.getPassword()));
        userRepository.save(user);
        log.info("User with ID {} updated successfully.", id);
    }

    @Override
    public void delete(UUID id) {
        log.info("Deleting user with ID: {}", id);
        if (userRepository.existsById(id)) {
            userRepository.deleteById(id);
            log.info("User with ID {} deleted successfully.", id);
        } else {
            log.warn("Attempted to delete non-existent user with ID: {}", id);
        }
    }
}
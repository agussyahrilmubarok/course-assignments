package com.example.user.service.impl;

import com.example.user.domain.User;
import com.example.user.exception.UserNotFoundException;
import com.example.user.model.UserDTO;
import com.example.user.repos.UserRepository;
import com.example.user.service.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Service
@Slf4j
@RequiredArgsConstructor
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;

    @Override
    public UserDTO findByID(String id) {
        User user = userRepository.findById(id)
                .orElseThrow(() -> {
                    log.error("User not found with id={}", id);
                    return new UserNotFoundException("User not found with id " + id);
                });

        return UserDTO.from(user);
    }
}

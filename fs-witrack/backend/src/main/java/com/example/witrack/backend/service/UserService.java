package com.example.witrack.backend.service;

import com.example.witrack.backend.model.UserDTO;

import java.util.UUID;

public interface UserService {

    UserDTO.UserResponse findById(UUID id);
}

package com.example.witrack.backend.service;

import com.example.witrack.backend.model.UserResponse;

public interface UserService {

    UserResponse getById(String id);
}

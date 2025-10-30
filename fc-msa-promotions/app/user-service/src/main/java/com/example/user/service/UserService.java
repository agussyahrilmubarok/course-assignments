package com.example.user.service;

import com.example.user.model.UserDTO;

public interface UserService {

    UserDTO findByID(final String id);
}

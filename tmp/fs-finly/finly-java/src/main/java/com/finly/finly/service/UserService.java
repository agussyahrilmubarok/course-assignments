package com.finly.finly.service;

import com.finly.finly.model.UserDTO;

import java.util.List;
import java.util.UUID;

public interface UserService {

    List<UserDTO> findAll();

    UserDTO get(final UUID id);

    UUID create(final UserDTO userDTO);

    void update(final UUID id, final UserDTO userDTO);

    void delete(final UUID id);
}

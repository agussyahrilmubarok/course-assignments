package com.finly.finly.model;

import com.finly.finly.domain.User;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;

import java.util.UUID;


@Getter
@Setter
public class UserDTO {

    private UUID id;

    @NotNull
    @Size(max = 255)
    @UserEmailUnique
    private String email;

    @NotNull
    @Size(max = 255)
    private String password;

    public static UserDTO fromUser(User user) {
        UserDTO userDTO = new UserDTO();
        user.setId(user.getId());
        user.setEmail(user.getEmail());
        user.setPassword(user.getPassword());
        return userDTO;
    }
}

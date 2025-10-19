package com.finly.finly.model;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@LoginPasswordMatchesRegistered(keyCheck = "email", valueCheck = "password")
public class LoginDTO {

    @NotBlank(message = "Email is required.")
    @Email(message = "Please provide a valid email address.")
    @LoginEmailMatchesRegistered(message = "Email is not registered.")
    private String email;

    @NotEmpty(message = "Password is required")
    @Size(min = 8, max = 255, message = "Password must be between 8 and 255 characters.")
    private String password;
}

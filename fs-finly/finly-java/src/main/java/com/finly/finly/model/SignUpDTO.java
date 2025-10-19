package com.finly.finly.model;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@SignUpPasswordMatches(sourceField = "password", targetField = "confirmPassword")
public class SignUpDTO {

    @NotBlank(message = "Email is required.")
    @Email(message = "Please provide a valid email address.")
    @SignUpEmailAvailable(message = "Email has been registered.")
    private String email;

    @NotBlank(message = "Password is required.")
    @Size(min = 8, max = 255, message = "Password must be between 8 and 255 characters.")
    private String password;

    @NotBlank(message = "Please confirm your password.")
    private String confirmPassword;
}

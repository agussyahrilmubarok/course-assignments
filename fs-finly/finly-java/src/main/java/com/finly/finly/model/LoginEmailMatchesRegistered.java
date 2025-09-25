package com.finly.finly.model;

import com.finly.finly.repos.UserRepository;
import jakarta.validation.Constraint;
import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;
import jakarta.validation.Payload;

import java.lang.annotation.*;

/**
 * Custom annotation to check if the email is already registered.
 */
@Target(ElementType.FIELD)
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Constraint(validatedBy = LoginEmailMatchesRegistered.EmailExistsValidator.class)
public @interface LoginEmailMatchesRegistered {

    String message() default "Email is already registered.";

    Class<?>[] groups() default {};

    Class<? extends Payload>[] payload() default {};

    // Validator for checking if the email exists in the system.
    class EmailExistsValidator implements ConstraintValidator<LoginEmailMatchesRegistered, String> {

        private final UserRepository userRepository;

        public EmailExistsValidator(UserRepository userRepository) {
            this.userRepository = userRepository;
        }

        @Override
        public boolean isValid(String email, ConstraintValidatorContext context) {
            // If email is null or empty, skip the validation.
            if (email == null || email.isEmpty()) {
                return true;
            }

            boolean isValid = true;
            boolean exists = userRepository.existsByEmailIgnoreCase(email);
            if (exists) return isValid;
            else return !isValid;
        }
    }
}

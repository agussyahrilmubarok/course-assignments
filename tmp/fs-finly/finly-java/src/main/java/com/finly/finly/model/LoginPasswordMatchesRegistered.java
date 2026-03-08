package com.finly.finly.model;

import com.finly.finly.domain.User;
import com.finly.finly.repos.UserRepository;
import jakarta.validation.Constraint;
import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;
import jakarta.validation.Payload;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.BeanWrapperImpl;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.lang.annotation.*;

/**
 * Custom annotation to validate if the password matches the registered password.
 */
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Constraint(validatedBy = LoginPasswordMatchesRegistered.PasswordMatchesValidator.class)
// Renamed to PasswordMatchesValidator
public @interface LoginPasswordMatchesRegistered {

    String message() default "Password is incorrect"; // More clear message

    Class<?>[] groups() default {};

    Class<? extends Payload>[] payload() default {};

    String keyCheck();

    String valueCheck();

    @Target({ElementType.TYPE})
    @Retention(RetentionPolicy.RUNTIME)
    @interface List {
        LoginPasswordMatchesRegistered[] value();
    }

    // Validator for password matching
    @RequiredArgsConstructor
    class PasswordMatchesValidator implements ConstraintValidator<LoginPasswordMatchesRegistered, Object> {

        private final UserRepository userRepository;
        private final PasswordEncoder passwordEncoder;

        private String keyCheck;
        private String valueCheck;
        private String message;

        @Override
        public void initialize(final LoginPasswordMatchesRegistered constraintAnnotation) {
            this.keyCheck = constraintAnnotation.keyCheck();
            this.valueCheck = constraintAnnotation.valueCheck();
            this.message = constraintAnnotation.message();
        }

        @Override
        public boolean isValid(final Object value, final ConstraintValidatorContext context) {

            // Retrieve field values dynamically based on annotations
            String fieldKey = (String) new BeanWrapperImpl(value).getPropertyValue(keyCheck);
            String fieldValue = (String) new BeanWrapperImpl(value).getPropertyValue(valueCheck);

            // If the email is null or empty, no need to validate
            if (fieldValue == null || fieldValue.equals("")) {
                return true;
            }

            // Find user by email
            User user = userRepository.findByEmail(fieldKey).orElse(null);
            if (user == null) {
                return false; // User not found, invalid password
            }

            // Check if password matches
            boolean isValid = passwordEncoder.matches(fieldValue, user.getPassword());
            if (!isValid) {
                context.disableDefaultConstraintViolation();
                context
                        .buildConstraintViolationWithTemplate(message)
                        .addPropertyNode(valueCheck)
                        .addConstraintViolation();
            }

            return isValid;
        }
    }
}

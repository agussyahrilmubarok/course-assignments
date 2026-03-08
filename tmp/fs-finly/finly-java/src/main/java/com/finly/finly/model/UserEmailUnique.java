package com.finly.finly.model;

import com.finly.finly.repos.UserRepository;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.validation.Constraint;
import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;
import jakarta.validation.Payload;
import org.springframework.web.servlet.HandlerMapping;

import java.lang.annotation.Documented;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;
import java.util.Map;
import java.util.UUID;

import static java.lang.annotation.ElementType.*;


/**
 * Validate that the email value isn't taken yet.
 */
@Target({FIELD, METHOD, ANNOTATION_TYPE})
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Constraint(
        validatedBy = UserEmailUnique.UserEmailUniqueValidator.class
)
public @interface UserEmailUnique {

    String message() default "{Exists.user.email}";

    Class<?>[] groups() default {};

    Class<? extends Payload>[] payload() default {};

    class UserEmailUniqueValidator implements ConstraintValidator<UserEmailUnique, String> {

        private final UserRepository userRepository;
        private final HttpServletRequest request;

        public UserEmailUniqueValidator(final UserRepository userRepository,
                                        final HttpServletRequest request) {
            this.userRepository = userRepository;
            this.request = request;
        }

        @Override
        public boolean isValid(final String value, final ConstraintValidatorContext cvContext) {
            if (value == null) {
                return true;
            }

            @SuppressWarnings("unchecked") final Map<String, String> pathVariables =
                    (Map<String, String>) request.getAttribute(HandlerMapping.URI_TEMPLATE_VARIABLES_ATTRIBUTE);

            final String currentId = pathVariables.get("id");

            if (currentId != null) {
                return userRepository.findById(UUID.fromString(currentId))
                        .map(user -> value.equalsIgnoreCase(user.getEmail()))
                        .orElse(!userRepository.existsByEmailIgnoreCase(value));
            }

            return !userRepository.existsByEmailIgnoreCase(value);
        }
    }
}

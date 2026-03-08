package com.finly.finly.model;

import com.finly.finly.repos.CustomerRepository;
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
        validatedBy = CustomerEmailUnique.CustomerEmailUniqueValidator.class
)
public @interface CustomerEmailUnique {

    String message() default "{Exists.customer.email}";

    Class<?>[] groups() default {};

    Class<? extends Payload>[] payload() default {};

    class CustomerEmailUniqueValidator implements ConstraintValidator<CustomerEmailUnique, String> {

        private final CustomerRepository customerRepository;
        private final HttpServletRequest request;

        public CustomerEmailUniqueValidator(final CustomerRepository customerRepository,
                                            final HttpServletRequest request) {
            this.customerRepository = customerRepository;
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
                return customerRepository.findById(UUID.fromString(currentId))
                        .map(customer -> value.equalsIgnoreCase(customer.getEmail()))
                        .orElse(!customerRepository.existsByEmailIgnoreCase(value));
            }

            return !customerRepository.existsByEmailIgnoreCase(value);
        }
    }
}

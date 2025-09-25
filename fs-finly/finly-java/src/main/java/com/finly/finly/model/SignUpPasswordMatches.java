package com.finly.finly.model;

import jakarta.validation.Constraint;
import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;
import jakarta.validation.Payload;
import org.springframework.beans.BeanWrapperImpl;

import java.lang.annotation.*;

/**
 * Custom annotation to validate that two password fields match.
 * Can be used to compare a "confirm password" field with a "password" field.
 */
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Constraint(validatedBy = SignUpPasswordMatches.PasswordMatcherValidator.class)
public @interface SignUpPasswordMatches {

    // The default error message when the passwords do not match.
    String message() default "Passwords do not match";

    // Groups can be used to specify different validation groups.
    Class<?>[] groups() default {};

    // Payload can be used to provide additional data to the validator.
    Class<? extends Payload>[] payload() default {};

    String sourceField();

    String targetField();


    @Target({ElementType.TYPE})
    @Retention(RetentionPolicy.RUNTIME)
    @interface List {
        SignUpPasswordMatches[] value();
    }

    // Validator that checks if the two password fields match
    class PasswordMatcherValidator implements ConstraintValidator<SignUpPasswordMatches, Object> {

        private String sourceField;
        private String targetField;
        private String message;

        @Override
        public void initialize(final SignUpPasswordMatches constraintAnnotation) {
            this.sourceField = constraintAnnotation.sourceField();
            this.targetField = constraintAnnotation.targetField();
            this.message = constraintAnnotation.message();
        }

        @Override
        public boolean isValid(final Object value, final ConstraintValidatorContext context) {
            Object fieldValue = new BeanWrapperImpl(value).getPropertyValue(sourceField);
            Object fieldMatchValue = new BeanWrapperImpl(value).getPropertyValue(targetField);

            boolean isValid = fieldValue != null && fieldValue.equals(fieldMatchValue);
            if (!isValid) {
                context.disableDefaultConstraintViolation();
                context
                        .buildConstraintViolationWithTemplate(message)
                        .addPropertyNode(targetField)
                        .addConstraintViolation();
            }

            return isValid;
        }
    }
}

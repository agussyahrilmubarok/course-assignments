package com.example.witrack.backend.model;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.repository.TicketRepository;
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

import static java.lang.annotation.ElementType.*;

@Target({FIELD, METHOD, ANNOTATION_TYPE})
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Constraint(
        validatedBy = TicketCodeUnique.TicketCodeUniqueValidator.class
)
public @interface TicketCodeUnique {

    String message() default "{Exists.ticket.code}";

    Class<?>[] groups() default {};

    Class<? extends Payload>[] payload() default {};

    class TicketCodeUniqueValidator implements ConstraintValidator<TicketCodeUnique, String> {

        private final TicketRepository ticketRepository;
        private final HttpServletRequest request;

        public TicketCodeUniqueValidator(final TicketRepository ticketRepository,
                                         final HttpServletRequest request) {
            this.ticketRepository = ticketRepository;
            this.request = request;
        }

        @Override
        public boolean isValid(final String value, final ConstraintValidatorContext cvContext) {
            if (value == null) {
                // no value present
                return true;
            }
            @SuppressWarnings("unchecked") final Map<String, String> pathVariables =
                    ((Map<String, String>) request.getAttribute(HandlerMapping.URI_TEMPLATE_VARIABLES_ATTRIBUTE));
            final String currentId = pathVariables.get("id");
            final Ticket ticket = ticketRepository.findById(currentId).orElseThrow();
            if (currentId != null && value.equalsIgnoreCase(ticket.getCode())) {
                // value hasn't changed
                return true;
            }
            return !ticketRepository.existsByCodeIgnoreCase(value);
        }

    }

}

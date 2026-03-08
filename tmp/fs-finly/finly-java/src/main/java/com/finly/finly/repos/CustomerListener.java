package com.finly.finly.repos;

import com.finly.finly.domain.Customer;
import org.springframework.data.mongodb.core.mapping.event.AbstractMongoEventListener;
import org.springframework.data.mongodb.core.mapping.event.BeforeConvertEvent;
import org.springframework.stereotype.Component;

import java.util.UUID;


@Component
public class CustomerListener extends AbstractMongoEventListener<Customer> {

    @Override
    public void onBeforeConvert(final BeforeConvertEvent<Customer> event) {
        if (event.getSource().getId() == null) {
            event.getSource().setId(UUID.randomUUID());
        }
    }

}

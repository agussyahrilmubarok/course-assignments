package com.finly.finly.repos;

import com.finly.finly.domain.User;
import org.springframework.data.mongodb.core.mapping.event.AbstractMongoEventListener;
import org.springframework.data.mongodb.core.mapping.event.BeforeConvertEvent;
import org.springframework.stereotype.Component;

import java.util.UUID;


@Component
public class UserListener extends AbstractMongoEventListener<User> {

    @Override
    public void onBeforeConvert(final BeforeConvertEvent<User> event) {
        if (event.getSource().getId() == null) {
            event.getSource().setId(UUID.randomUUID());
        }
    }

}

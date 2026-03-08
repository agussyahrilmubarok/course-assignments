package com.finly.finly.repos;

import com.finly.finly.domain.User;
import org.springframework.data.mongodb.repository.MongoRepository;

import java.util.Optional;
import java.util.UUID;


public interface UserRepository extends MongoRepository<User, UUID> {

    Optional<User> findByEmail(String email);

    boolean existsByEmailIgnoreCase(String email);
}

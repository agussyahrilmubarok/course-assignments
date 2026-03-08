package com.example.witrack.backend.repository;

import com.example.witrack.backend.domain.User;
import org.springframework.data.mongodb.repository.MongoRepository;

import java.util.Optional;

public interface UserRepository extends MongoRepository<User, String> {

    Optional<User> findByEmail(String email);

    boolean existsByIdIgnoreCase(String id);

    boolean existsByEmailIgnoreCase(String email);
}

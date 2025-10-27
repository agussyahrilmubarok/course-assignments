package com.example.grpc.server.repos;

import com.example.grpc.server.domain.AuthorEntity;
import org.springframework.data.jpa.repository.JpaRepository;


public interface AuthorRepository extends JpaRepository<AuthorEntity, Long> {
}

package com.example.grpc.server.repos;

import com.example.grpc.server.domain.BookEntity;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;


public interface BookRepository extends JpaRepository<BookEntity, Long> {

    List<BookEntity> findAllByAuthors_Id(Long id);

}

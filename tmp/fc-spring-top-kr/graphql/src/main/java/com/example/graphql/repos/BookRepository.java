package com.example.graphql.repos;

import com.example.graphql.domain.Book;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;


public interface BookRepository extends JpaRepository<Book, Long> {

    List<Book> findAllByAuthorsId(Long id);

}

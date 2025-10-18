package com.example.rest.repos;

import com.example.rest.domain.Review;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;


public interface ReviewRepository extends JpaRepository<Review, Long> {

    List<Review> findAllByBookId(Long bookId);

    Review findFirstByBookId(Long id);
}

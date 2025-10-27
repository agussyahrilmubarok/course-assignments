package com.example.grpc.server.repos;

import com.example.grpc.server.domain.ReviewEntity;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;


public interface ReviewRepository extends JpaRepository<ReviewEntity, Long> {

    List<ReviewEntity> findAllByBook_Id(Long bookId);

    ReviewEntity findFirstByBook_Id(Long id);
}

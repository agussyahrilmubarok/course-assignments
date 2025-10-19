package com.example.point.repos;

import com.example.point.domain.Point;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;

public interface PointRepository extends JpaRepository<Point, String> {

    List<Point> findByUserIdOrderByCreatedAtDesc(String userId);

    @Query("SELECT p FROM Point p " +
            "LEFT JOIN FETCH p.pointBalance " +
            "WHERE p.userId = :userId " +
            "ORDER BY p.createdAt DESC")
    Page<Point> findByUserIdOrderByCreatedAtDesc(@Param("userId") String userId, Pageable pageable);
}

package com.example.timesale.repos;

import com.example.timesale.domain.TimeSale;
import jakarta.persistence.LockModeType;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.time.LocalDateTime;
import java.util.Optional;

public interface TimeSaleRepository extends JpaRepository<TimeSale, String> {

    @EntityGraph(attributePaths = {"product"})
    @Query("SELECT ts FROM TimeSale ts " +
            "WHERE ts.startAt <= :now AND ts.endAt > :now AND ts.status = :status")
    Page<TimeSale> findAllByStartAtBeforeAndEndAtAfterAndStatus(@Param("now") LocalDateTime now,
                                                                @Param("status") TimeSale.Status status,
                                                                Pageable pageable);

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @EntityGraph(attributePaths = "product")
    @Query("SELECT ts FROM TimeSale ts WHERE ts.id = :id")
    Optional<TimeSale> findByIdWithPessimisticLock(@Param("id") String id);

    @Override
    @EntityGraph(attributePaths = "product")
    Optional<TimeSale> findById(String id);
}

package com.example.point.batch.point.repos;

import com.example.point.batch.point.domain.PointBalance;
import jakarta.persistence.LockModeType;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Lock;

import java.util.Optional;

public interface PointBalanceRepository extends JpaRepository<PointBalance, String> {

    @Lock(LockModeType.OPTIMISTIC)
    Optional<PointBalance> findByUserId(String userId);
}
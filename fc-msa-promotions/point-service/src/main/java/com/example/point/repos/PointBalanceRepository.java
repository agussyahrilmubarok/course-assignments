package com.example.point.repos;

import com.example.point.domain.PointBalance;
import org.springframework.data.jpa.repository.JpaRepository;

public interface PointBalanceRepository extends JpaRepository<PointBalance, String> {
}

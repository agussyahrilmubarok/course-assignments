package com.example.coupon.repos;

import com.example.coupon.domain.CouponPolicy;
import jakarta.persistence.LockModeType;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Query;

import java.util.Optional;

public interface CouponPolicyRepository extends JpaRepository<CouponPolicy, String> {

    /**
     * The reason for using PESSIMISTIC_WRITE is to ensure data consistency.
     * It is used to prevent conflicts and maintain data integrity when multiple transactions attempt
     * to modify the same data simultaneously.
     */
    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("SELECT cp FROM CouponPolicy cp WHERE cp.id = :id")
    Optional<CouponPolicy> findByIdWithLock(String id);
}

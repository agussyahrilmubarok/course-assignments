package com.example.coupon.repos;

import com.example.coupon.domain.Coupon;
import com.example.coupon.domain.CouponPolicy;
import jakarta.persistence.LockModeType;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.Optional;

public interface CouponRepository extends JpaRepository<Coupon, String> {

    /**
     * The reason for using PESSIMISTIC_WRITE is to ensure data consistency.
     * It is used to prevent conflicts and maintain data integrity when multiple transactions attempt
     * to modify the same data simultaneously.
     */
    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("SELECT c FROM Coupon c WHERE c.id = :id")
    Optional<Coupon> findByIdWithLock(@Param("id") String id);

    Page<Coupon> findByUserIdAndStatusOrderByCreatedAtDesc(String userId, Coupon.Status status, Pageable pageable);

    @EntityGraph(attributePaths = {"couponPolicy"})
    Optional<Coupon> findByIdAndUserId(String id, String userId);

    Coupon findFirstByCouponPolicy(CouponPolicy couponPolicy);

    @Query("SELECT COUNT(c) FROM Coupon c WHERE c.couponPolicy.id = :policyId")
    Long countByCouponPolicyId(@Param("policyId") String policyId);

    boolean existsByCodeIgnoreCase(String code);

}

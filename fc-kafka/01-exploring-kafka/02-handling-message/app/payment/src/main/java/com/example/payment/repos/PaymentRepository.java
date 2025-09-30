package com.example.payment.repos;

import com.example.payment.domain.Payment;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;


public interface PaymentRepository extends JpaRepository<Payment, UUID> {
}

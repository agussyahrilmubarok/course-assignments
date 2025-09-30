package com.example.payment.repos;

import com.example.payment.domain.Payment;
import java.util.UUID;
import org.springframework.data.jpa.repository.JpaRepository;


public interface PaymentRepository extends JpaRepository<Payment, UUID> {
}

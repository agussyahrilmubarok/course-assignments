package com.example.payment.repos;

import com.example.payment.domain.Transaction;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;

public interface TransactionRepository extends JpaRepository<Transaction, String> {

    Optional<Transaction> findByIdAndOrderId(String id, String orderId);
}

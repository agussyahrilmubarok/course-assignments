package com.example.order.repos;

import com.example.order.domain.Order;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;


public interface OrderRepository extends JpaRepository<Order, UUID> {
}

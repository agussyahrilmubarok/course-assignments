package com.example.order.repos;

import com.example.order.domain.Order;
import java.util.UUID;
import org.springframework.data.jpa.repository.JpaRepository;


public interface OrderRepository extends JpaRepository<Order, UUID> {
}

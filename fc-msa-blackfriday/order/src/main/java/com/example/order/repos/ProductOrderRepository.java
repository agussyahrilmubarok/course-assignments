package com.example.order.repos;

import com.example.order.domain.ProductOrder;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface ProductOrderRepository extends JpaRepository<ProductOrder, String> {

    List<ProductOrder> findAllByUserId(String userId);
}

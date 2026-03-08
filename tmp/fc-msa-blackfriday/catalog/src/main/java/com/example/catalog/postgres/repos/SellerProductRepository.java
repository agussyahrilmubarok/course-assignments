package com.example.catalog.postgres.repos;

import com.example.catalog.postgres.domain.SellerProduct;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Optional;

public interface SellerProductRepository extends JpaRepository<SellerProduct, String> {

    List<SellerProduct> findBySellerId(String sellerId);

    Optional<SellerProduct> findByProductId(String productId);

    void deleteByProductId(String productId);
}

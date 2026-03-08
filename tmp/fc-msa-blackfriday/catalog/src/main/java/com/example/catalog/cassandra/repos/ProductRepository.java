package com.example.catalog.cassandra.repos;

import com.example.catalog.cassandra.domain.Product;
import org.springframework.data.cassandra.repository.CassandraRepository;
import org.springframework.data.cassandra.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface ProductRepository extends CassandraRepository<Product, String> {

    //TODO: change seller id to indexing
    @Query("SELECT * FROM products WHERE sellerId = ?0 ALLOW FILTERING")
    List<Product> findAllBySellerId(String sellerId);

    @Query("SELECT * FROM products WHERE tags CONTAINS ?0 ALLOW FILTERING")
    List<Product> findAllByTag(String tag);

    @Query("SELECT * FROM products WHERE price >= ?0 AND price <= ?1 ALLOW FILTERING")
    List<Product> findAllByPriceRange(Long minPrice, Long maxPrice);

}

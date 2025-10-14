package com.example.catalog.cassandra.repos;

import com.example.catalog.cassandra.domain.Product;
import org.springframework.data.cassandra.repository.CassandraRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface ProductRepository extends CassandraRepository<Product, String> {
}

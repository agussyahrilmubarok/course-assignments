package com.finly.finly.repos;

import com.finly.finly.domain.Customer;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.data.mongodb.repository.Query;

import java.util.List;
import java.util.UUID;


public interface CustomerRepository extends MongoRepository<Customer, UUID> {

    @Query("{ 'user': ?0, $or: [ " +
            "{ 'name': { $regex: ?1, $options: 'i' } }, " +
            "{ 'email': { $regex: ?1, $options: 'i' } }, " +
            "{ 'phone': { $regex: ?1, $options: 'i' } }, " +
            "{ 'address': { $regex: ?1, $options: 'i' } } " +
            "] }")
    List<Customer> searchCustomersByUserId(UUID userId, String search);

    List<Customer> findByUserId(UUID userId);

    Customer findFirstByUserId(UUID id);

    boolean existsByEmailIgnoreCase(String email);

    long countByUserId(UUID userId);
}

package com.finly.finly.repos;

import com.finly.finly.domain.Customer;

import java.util.List;
import java.util.UUID;

public interface CustomerCustomRepository {

    List<Customer> findCustomersByUserAndSearch(UUID userId, String search);

    List<Customer> findCustomersByUserId(UUID userId);
}

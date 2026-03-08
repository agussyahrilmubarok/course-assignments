package com.finly.finly.service;

import com.finly.finly.model.CustomerDTO;

import java.util.List;
import java.util.UUID;

public interface CustomerService {

    List<CustomerDTO> findAll();

    List<CustomerDTO> findCustomersByOwnerAndSearch(UUID userId, String search);

    List<CustomerDTO> findByUserId(UUID userId);

    CustomerDTO get(UUID id);

    UUID create(CustomerDTO customerDTO);

    void update(UUID id, CustomerDTO customerDTO);

    void delete(UUID id);

    long countByUser(UUID userId);
}

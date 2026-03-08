package com.finly.finly.service.impl;

import com.finly.finly.domain.Customer;
import com.finly.finly.domain.User;
import com.finly.finly.model.CustomerDTO;
import com.finly.finly.repos.CustomerCustomRepository;
import com.finly.finly.repos.CustomerRepository;
import com.finly.finly.repos.UserRepository;
import com.finly.finly.service.CustomerService;
import com.finly.finly.util.NotFoundException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class CustomerServiceImpl implements CustomerService {

    private final CustomerRepository customerRepository;
    private final CustomerCustomRepository customerCustomRepository;
    private final UserRepository userRepository;

    @Override
    public List<CustomerDTO> findAll() {
        log.info("Retrieving all customers from the database.");
        return customerRepository.findAll().stream()
                .map(CustomerDTO::fromCustomer)
                .toList();
    }

    @Override
    public List<CustomerDTO> findCustomersByOwnerAndSearch(UUID userId, String search) {
        log.info("Retrieving customers for user ID {} with search query: {}", userId, search);
        List<Customer> customers = (search != null && !search.isEmpty())
                ? customerCustomRepository.findCustomersByUserAndSearch(userId, search)
                : customerCustomRepository.findCustomersByUserId(userId);
        return customers.stream()
                .map(CustomerDTO::fromCustomer)
                .toList();
    }

    @Override
    public List<CustomerDTO> findByUserId(UUID userId) {
        log.info("Retrieving all customers associated with user ID: {}", userId);
        return customerRepository.findByUserId(userId).stream()
                .map(CustomerDTO::fromCustomer)
                .toList();
    }

    @Override
    public CustomerDTO get(UUID id) {
        log.info("Retrieving customer details for ID: {}", id);
        return customerRepository.findById(id)
                .map(CustomerDTO::fromCustomer)
                .orElseThrow(() -> {
                    log.error("Customer with ID {} not found.", id);
                    return new NotFoundException("Customer not found for ID: " + id);
                });
    }

    @Override
    public UUID create(CustomerDTO customerDTO) {
        log.info("Creating a new customer for user ID: {}", customerDTO.getId());
        User user = userRepository.findById(customerDTO.getUser())
                .orElseThrow(() -> {
                    log.error("User with ID {} not found. Customer creation aborted.", customerDTO.getUser());
                    return new NotFoundException("User not found for ID: " + customerDTO.getUser());
                });

        Customer customer = new Customer();
        customer.setName(customerDTO.getName());
        customer.setEmail(customerDTO.getEmail());
        customer.setPhone(customerDTO.getPhone());
        customer.setAddress(customerDTO.getAddress());
        customer.setUser(user);
        UUID createdId = customerRepository.save(customer).getId();
        log.info("Customer created successfully with ID: {}", createdId);
        return createdId;
    }

    @Override
    public void update(UUID id, CustomerDTO customerDTO) {
        log.info("Updating customer with ID: {}", id);
        Customer customer = customerRepository.findById(id)
                .orElseThrow(() -> {
                    log.error("Customer with ID {} not found. Update aborted.", id);
                    return new NotFoundException("Customer not found for ID: " + id);
                });

        User user = userRepository.findById(customerDTO.getUser())
                .orElseThrow(() -> {
                    log.error("User with ID {} not found. Update aborted for customer ID: {}", customerDTO.getUser(), id);
                    return new NotFoundException("User not found for ID: " + customerDTO.getUser());
                });

        customer.setName(customerDTO.getName());
        customer.setEmail(customerDTO.getEmail());
        customer.setPhone(customerDTO.getPhone());
        customer.setAddress(customerDTO.getAddress());
        customer.setUser(user);
        customerRepository.save(customer);
        log.info("Customer with ID {} updated successfully.", id);
    }

    @Override
    public void delete(UUID id) {
        log.info("Deleting customer with ID: {}", id);
        if (customerRepository.existsById(id)) {
            customerRepository.deleteById(id);
            log.info("Customer with ID {} deleted successfully.", id);
        } else {
            log.warn("Attempted to delete non-existent customer with ID: {}", id);
        }
    }

    @Override
    public long countByUser(UUID userId) {
        log.info("Counting customers for user ID: {}", userId);
        return customerRepository.countByUserId(userId);
    }
}
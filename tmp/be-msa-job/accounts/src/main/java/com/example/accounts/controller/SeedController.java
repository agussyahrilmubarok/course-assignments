package com.example.accounts.controller;

import com.example.accounts.entity.Accounts;
import com.example.accounts.entity.Customer;
import com.example.accounts.repository.AccountsRepository;
import com.example.accounts.repository.CustomerRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/seeds")
@RequiredArgsConstructor
public class SeedController {

    private final CustomerRepository customerRepository;
    private final AccountsRepository accountsRepository;

    @GetMapping("/init")
    public ResponseEntity<List<Accounts>> initializeData() {
        List<Accounts> accounts = List.of();
        if (customerRepository.count() > 0) {
            return ResponseEntity.ok(accounts);
        }

        // Insert customers one-by-one to capture the generated ID
        Customer c1 = customerRepository.save(new Customer(null, "John Doe", "john.doe@example.com", "0812345678"));

        // Insert accounts with actual customer IDs
        Accounts a1 = accountsRepository.save(new Accounts(c1.getCustomerId(), 1000001L, "Savings", "Jakarta"));

        accounts = List.of(a1);

        return ResponseEntity.ok(accounts);
    }

    @GetMapping("/clear")
    public String clearData() {
        accountsRepository.deleteAll();
        customerRepository.deleteAll();
        return "🗑️ All customers and accounts deleted successfully.";
    }
}


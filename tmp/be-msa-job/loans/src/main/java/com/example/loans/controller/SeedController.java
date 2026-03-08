package com.example.loans.controller;

import com.example.loans.entity.Loans;
import com.example.loans.repository.LoansRepository;
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

    private final LoansRepository loansRepository;

    @GetMapping("/init")
    public ResponseEntity<List<Loans>> initializeLoans() {
        List<Loans> loans = List.of();

        // Check if data already exists
        if (loansRepository.count() > 0) {
            return ResponseEntity.ok(loans);
        }

        // Create loans for customers
        Loans loan1 = new Loans(null, "0812345678", "LN10001", "Personal Loan", 5000000, 1000000, 4000000);

        // Save all loans
        loans = loansRepository.saveAll(List.of(loan1));

        return ResponseEntity.ok(loans);

    }

    @GetMapping("/clear")
    public String clearLoans() {
        loansRepository.deleteAll();
        return "✅ All loan records cleared.";
    }
}

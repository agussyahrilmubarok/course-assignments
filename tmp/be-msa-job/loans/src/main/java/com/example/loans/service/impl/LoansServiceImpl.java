package com.example.loans.service.impl;

import com.example.loans.constants.LoansConstants;
import com.example.loans.dto.LoansDto;
import com.example.loans.entity.Loans;
import com.example.loans.exception.LoanAlreadyExistsException;
import com.example.loans.exception.ResourceNotFoundException;
import com.example.loans.mapper.LoansMapper;
import com.example.loans.repository.LoansRepository;
import com.example.loans.service.ILoansService;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.Optional;
import java.util.Random;

@Service
@Slf4j
@AllArgsConstructor
public class LoansServiceImpl implements ILoansService {

    private LoansRepository loansRepository;

    /**
     * @param mobileNumber - Mobile Number of the Customer
     */
    @Override
    public void createLoan(String mobileNumber) {
        log.info("Attempting to create loan for mobileNumber: {}", mobileNumber);
        Optional<Loans> optionalLoans = loansRepository.findByMobileNumber(mobileNumber);
        if (optionalLoans.isPresent()) {
            log.error("Loan already exists for mobileNumber: {}", mobileNumber);
            throw new LoanAlreadyExistsException("Loan already registered with given mobileNumber " + mobileNumber);
        }

        Loans newLoan = createNewLoan(mobileNumber);
        loansRepository.save(newLoan);
        log.info("Loan successfully created with loanNumber: {} for mobileNumber: {}", newLoan.getLoanNumber(), mobileNumber);
    }

    /**
     * @param mobileNumber - Mobile Number of the Customer
     * @return the new loan details
     */
    private Loans createNewLoan(String mobileNumber) {
        long randomLoanNumber = 100000000000L + new Random().nextInt(900000000);
        Loans newLoan = new Loans();
        newLoan.setLoanNumber(Long.toString(randomLoanNumber));
        newLoan.setMobileNumber(mobileNumber);
        newLoan.setLoanType(LoansConstants.HOME_LOAN);
        newLoan.setTotalLoan(LoansConstants.NEW_LOAN_LIMIT);
        newLoan.setAmountPaid(0);
        newLoan.setOutstandingAmount(LoansConstants.NEW_LOAN_LIMIT);
        log.debug("Initialized new loan object for mobileNumber: {}, loanNumber: {}", mobileNumber, randomLoanNumber);
        return newLoan;
    }

    /**
     * @param mobileNumber - Input mobile Number
     * @return Loan Details based on a given mobileNumber
     */
    @Override
    public LoansDto fetchLoan(String mobileNumber) {
        log.info("Fetching loan details for mobileNumber: {}", mobileNumber);
        Loans loans = loansRepository.findByMobileNumber(mobileNumber).orElseThrow(() -> {
            log.error("Loan not found for mobileNumber: {}", mobileNumber);
            return new ResourceNotFoundException("Loan", "mobileNumber", mobileNumber);
        });

        LoansDto loansDto = LoansMapper.mapToLoansDto(loans, new LoansDto());
        log.info("Successfully fetched loan details for mobileNumber: {}", mobileNumber);
        return loansDto;
    }

    /**
     * @param loansDto - LoansDto Object
     * @return boolean indicating if the update of loan details is successful or not
     */
    @Override
    public boolean updateLoan(LoansDto loansDto) {
        log.info("Updating loan with loanNumber: {}", loansDto.getLoanNumber());
        Loans loans = loansRepository.findByLoanNumber(loansDto.getLoanNumber()).orElseThrow(() -> {
            log.error("Loan not found for loanNumber: {}", loansDto.getLoanNumber());
            return new ResourceNotFoundException("Loan", "LoanNumber", loansDto.getLoanNumber());
        });

        LoansMapper.mapToLoans(loansDto, loans);
        loansRepository.save(loans);
        log.info("Loan updated successfully for loanNumber: {}", loansDto.getLoanNumber());
        return true;
    }

    /**
     * @param mobileNumber - Input MobileNumber
     * @return boolean indicating if the delete of loan details is successful or not
     */
    @Override
    public boolean deleteLoan(String mobileNumber) {
        log.info("Deleting loan for mobileNumber: {}", mobileNumber);
        Loans loans = loansRepository.findByMobileNumber(mobileNumber).orElseThrow(() -> {
            log.error("Loan not found for mobileNumber: {}", mobileNumber);
            return new ResourceNotFoundException("Loan", "mobileNumber", mobileNumber);
        });

        loansRepository.deleteById(loans.getLoanId());
        log.info("Loan deleted successfully for mobileNumber: {}", mobileNumber);
        return true;
    }
}

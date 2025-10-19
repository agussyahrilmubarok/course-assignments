package com.example.accounts.service.impl;

import com.example.accounts.dto.AccountsDto;
import com.example.accounts.dto.CardsDto;
import com.example.accounts.dto.CustomerDetailsDto;
import com.example.accounts.dto.LoansDto;
import com.example.accounts.entity.Accounts;
import com.example.accounts.entity.Customer;
import com.example.accounts.exception.ResourceNotFoundException;
import com.example.accounts.mapper.AccountsMapper;
import com.example.accounts.mapper.CustomerMapper;
import com.example.accounts.repository.AccountsRepository;
import com.example.accounts.repository.CustomerRepository;
import com.example.accounts.service.ICustomersService;
import com.example.accounts.service.client.CardsFeignClient;
import com.example.accounts.service.client.LoansFeignClient;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;

@Service
@Slf4j
@AllArgsConstructor
public class CustomersServiceImpl implements ICustomersService {

    private AccountsRepository accountsRepository;
    private CustomerRepository customerRepository;
    private CardsFeignClient cardsFeignClient;
    private LoansFeignClient loansFeignClient;

    /**
     * @param mobileNumber  - Input Mobile Number
     * @param correlationId - Correlation ID value generated at Edge server
     * @return Customer Details based on a given mobileNumber
     */
    @Override
    public CustomerDetailsDto fetchCustomerDetails(String correlationId, String mobileNumber) {
        log.info("Fetching customer details for mobileNumber: {}, correlationId: {}", mobileNumber, correlationId);

        Customer customer = customerRepository.findByMobileNumber(mobileNumber).orElseThrow(() -> {
            log.error("Customer not found for mobileNumber: {}", mobileNumber);
            return new ResourceNotFoundException("Customer", "mobileNumber", mobileNumber);
        });

        Accounts accounts = accountsRepository.findByCustomerId(customer.getCustomerId()).orElseThrow(() -> {
            log.error("Accounts not found for customerId: {}", customer.getCustomerId());
            return new ResourceNotFoundException("Account", "customerId", customer.getCustomerId().toString());
        });

        CustomerDetailsDto customerDetailsDto = CustomerMapper.mapToCustomerDetailsDto(customer, new CustomerDetailsDto());
        customerDetailsDto.setAccountsDto(AccountsMapper.mapToAccountsDto(accounts, new AccountsDto()));

        log.info("Customer and accounts data fetched successfully for mobileNumber: {}", mobileNumber);

        try {
            ResponseEntity<LoansDto> loansDtoResponseEntity = loansFeignClient.fetchLoanDetails(correlationId, mobileNumber);
            if (loansDtoResponseEntity != null && loansDtoResponseEntity.getBody() != null) {
                customerDetailsDto.setLoansDto(loansDtoResponseEntity.getBody());
                log.info("Loan details fetched for mobileNumber: {}", mobileNumber);
            } else {
                log.warn("No loan details found for mobileNumber: {}", mobileNumber);
            }
        } catch (Exception e) {
            log.error("Error fetching loan details for mobileNumber: {} - {}", mobileNumber, e.getMessage());
        }

        try {
            CardsDto cardsDtoEntity = cardsFeignClient.fetchCardDetails(correlationId, mobileNumber);
            if (cardsDtoEntity != null) {
                customerDetailsDto.setCardsDto(cardsDtoEntity);
                log.info("Card details fetched for mobileNumber: {}", mobileNumber);
            } else {
                log.warn("No card details found for mobileNumber: {}", mobileNumber);
            }
        } catch (Exception e) {
            log.error("Error fetching card details for mobileNumber: {} - {}", mobileNumber, e.getMessage());
        }

        log.info("Completed fetching customer details for mobileNumber: {}", mobileNumber);
        return customerDetailsDto;
    }
}

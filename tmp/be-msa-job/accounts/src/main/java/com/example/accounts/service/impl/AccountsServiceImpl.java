package com.example.accounts.service.impl;

import com.example.accounts.constants.AccountsConstants;
import com.example.accounts.dto.AccountsDto;
import com.example.accounts.dto.CustomerDto;
import com.example.accounts.entity.Accounts;
import com.example.accounts.entity.Customer;
import com.example.accounts.exception.CustomerAlreadyExistsException;
import com.example.accounts.exception.ResourceNotFoundException;
import com.example.accounts.mapper.AccountsMapper;
import com.example.accounts.mapper.CustomerMapper;
import com.example.accounts.repository.AccountsRepository;
import com.example.accounts.repository.CustomerRepository;
import com.example.accounts.service.IAccountsService;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.Optional;
import java.util.Random;

@Service
@Slf4j
@AllArgsConstructor
public class AccountsServiceImpl implements IAccountsService {

    private AccountsRepository accountsRepository;
    private CustomerRepository customerRepository;

    /**
     * @param customerDto - CustomerDto Object
     */
    @Override
    public void createAccount(CustomerDto customerDto) {
        log.info("Creating account for mobile number: {}", customerDto.getMobileNumber());

        Customer customer = CustomerMapper.mapToCustomer(customerDto, new Customer());
        Optional<Customer> optionalCustomer = customerRepository.findByMobileNumber(customerDto.getMobileNumber());

        if (optionalCustomer.isPresent()) {
            log.warn("Customer already exists with mobile number: {}", customerDto.getMobileNumber());
            throw new CustomerAlreadyExistsException("Customer already registered with given mobileNumber "
                    + customerDto.getMobileNumber());
        }

        Customer savedCustomer = customerRepository.save(customer);
        log.info("Customer saved with ID: {}", savedCustomer.getCustomerId());

        Accounts account = accountsRepository.save(createNewAccount(savedCustomer));
        log.info("Account created with Account Number: {}", account.getAccountNumber());
    }

    /**
     * @param customer - Customer Object
     * @return the new account details
     */
    private Accounts createNewAccount(Customer customer) {
        Accounts newAccount = new Accounts();
        newAccount.setCustomerId(customer.getCustomerId());
        long randomAccNumber = 1000000000L + new Random().nextInt(900000000);
        newAccount.setAccountNumber(randomAccNumber);
        newAccount.setAccountType(AccountsConstants.SAVINGS);
        newAccount.setBranchAddress(AccountsConstants.ADDRESS);

        log.debug("Generated new account for customer ID {}: {}", customer.getCustomerId(), randomAccNumber);

        return newAccount;
    }

    /**
     * @param mobileNumber - Input Mobile Number
     * @return Accounts Details based on a given mobileNumber
     */
    @Override
    public CustomerDto fetchAccount(String mobileNumber) {
        log.info("Fetching account for mobile number: {}", mobileNumber);

        Customer customer = customerRepository.findByMobileNumber(mobileNumber).orElseThrow(
                () -> {
                    log.error("Customer not found with mobile number: {}", mobileNumber);
                    return new ResourceNotFoundException("Customer", "mobileNumber", mobileNumber);
                }
        );

        Accounts accounts = accountsRepository.findByCustomerId(customer.getCustomerId()).orElseThrow(
                () -> {
                    log.error("Account not found for customer ID: {}", customer.getCustomerId());
                    return new ResourceNotFoundException("Account", "customerId", customer.getCustomerId().toString());
                }
        );

        log.info("Fetched account number {} for customer ID {}", accounts.getAccountNumber(), customer.getCustomerId());

        CustomerDto customerDto = CustomerMapper.mapToCustomerDto(customer, new CustomerDto());
        customerDto.setAccountsDto(AccountsMapper.mapToAccountsDto(accounts, new AccountsDto()));
        return customerDto;
    }

    /**
     * @param customerDto - CustomerDto Object
     * @return boolean indicating if the update of Account details is successful or not
     */
    public boolean updateAccount(CustomerDto customerDto) {
        log.info("Updating account for mobile number: {}", customerDto.getMobileNumber());

        boolean isUpdated = false;
        AccountsDto accountsDto = customerDto.getAccountsDto();

        if (accountsDto != null) {
            Accounts accounts = accountsRepository.findById(accountsDto.getAccountNumber()).orElseThrow(
                    () -> {
                        log.error("Account not found with account number: {}", accountsDto.getAccountNumber());
                        return new ResourceNotFoundException("Account", "AccountNumber", accountsDto.getAccountNumber().toString());
                    }
            );

            AccountsMapper.mapToAccounts(accountsDto, accounts);
            accounts = accountsRepository.save(accounts);
            log.info("Account updated with account number: {}", accounts.getAccountNumber());

            Long customerId = accounts.getCustomerId();
            Customer customer = customerRepository.findById(customerId).orElseThrow(
                    () -> {
                        log.error("Customer not found with ID: {}", customerId);
                        return new ResourceNotFoundException("Customer", "CustomerID", customerId.toString());
                    }
            );

            CustomerMapper.mapToCustomer(customerDto, customer);
            customerRepository.save(customer);
            log.info("Customer updated with ID: {}", customerId);

            isUpdated = true;
        }

        return isUpdated;
    }

    /**
     * @param mobileNumber - Input Mobile Number
     * @return boolean indicating if the delete of Account details is successful or not
     */
    @Override
    public boolean deleteAccount(String mobileNumber) {
        log.info("Deleting account for mobile number: {}", mobileNumber);

        Customer customer = customerRepository.findByMobileNumber(mobileNumber).orElseThrow(
                () -> {
                    log.error("Customer not found with mobile number: {}", mobileNumber);
                    return new ResourceNotFoundException("Customer", "mobileNumber", mobileNumber);
                }
        );

        accountsRepository.deleteByCustomerId(customer.getCustomerId());
        log.info("Deleted account(s) for customer ID: {}", customer.getCustomerId());

        customerRepository.deleteById(customer.getCustomerId());
        log.info("Deleted customer with ID: {}", customer.getCustomerId());

        return true;
    }
}

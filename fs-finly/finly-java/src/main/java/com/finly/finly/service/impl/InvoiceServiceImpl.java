package com.finly.finly.service.impl;

import com.finly.finly.domain.Customer;
import com.finly.finly.domain.Invoice;
import com.finly.finly.domain.User;
import com.finly.finly.model.InvoiceDTO;
import com.finly.finly.model.InvoiceDetailDTO;
import com.finly.finly.repos.CustomerRepository;
import com.finly.finly.repos.InvoiceCustomRepository;
import com.finly.finly.repos.InvoiceRepository;
import com.finly.finly.repos.UserRepository;
import com.finly.finly.service.InvoiceService;
import com.finly.finly.util.NotFoundException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class InvoiceServiceImpl implements InvoiceService {

    private final InvoiceRepository invoiceRepository;
    private final InvoiceCustomRepository invoiceCustomRepository;
    private final CustomerRepository customerRepository;
    private final UserRepository userRepository;

    @Override
    public List<InvoiceDTO> findAll() {
        log.info("Retrieving all invoices from the database.");
        return invoiceRepository.findAll()
                .stream()
                .map(InvoiceDTO::fromInvoice)
                .toList();
    }

    @Override
    public List<InvoiceDetailDTO> findInvoiceByOwnerAndSearch(UUID userId, String search) {
        log.info("Retrieving invoices for user ID {}. Search term: '{}'", userId, search != null ? search : "none");
        List<Invoice> invoices = (search != null && !search.isEmpty())
                ? invoiceCustomRepository.findInvoicesByOwnerAndSearch(userId, search)
                : invoiceCustomRepository.findInvoicesByOwnerId(userId);
        return invoices.stream()
                .map(InvoiceDetailDTO::fromInvoice)
                .toList();
    }

    @Override
    public InvoiceDTO get(UUID id) {
        log.info("Retrieving invoice details for ID: {}", id);
        return invoiceRepository.findById(id)
                .map(InvoiceDTO::fromInvoice)
                .orElseThrow(() -> {
                    log.error("Invoice with ID {} not found.", id);
                    return new NotFoundException("Invoice not found for ID: " + id);
                });
    }

    @Override
    public InvoiceDetailDTO getDetail(UUID id) {
        log.info("Retrieving detailed invoice information for ID: {}", id);
        return invoiceRepository.findById(id)
                .map(InvoiceDetailDTO::fromInvoice)
                .orElseThrow(() -> {
                    log.error("Invoice with ID {} not found.", id);
                    return new NotFoundException("Invoice not found for ID: " + id);
                });
    }

    @Override
    public UUID create(InvoiceDTO invoiceDTO) {
        log.info("Creating new invoice for owner ID: {} and customer ID: {}", invoiceDTO.getOwner(), invoiceDTO.getCustomer());

        User owner = userRepository.findById(invoiceDTO.getOwner())
                .orElseThrow(() -> {
                    log.error("User with ID {} not found. Invoice creation aborted.", invoiceDTO.getOwner());
                    return new NotFoundException("User not found for ID: " + invoiceDTO.getOwner());
                });

        Customer customer = customerRepository.findById(invoiceDTO.getCustomer())
                .orElseThrow(() -> {
                    log.error("Customer with ID {} not found. Invoice creation aborted.", invoiceDTO.getCustomer());
                    return new NotFoundException("Customer not found for ID: " + invoiceDTO.getCustomer());
                });

        Invoice invoice = new Invoice();
        invoice.setAmount(invoiceDTO.getAmount());
        invoice.setDueDate(invoiceDTO.getDueDate());
        invoice.setStatus(invoiceDTO.getStatus());
        invoice.setOwner(owner);
        invoice.setCustomer(customer);

        UUID createdId = invoiceRepository.save(invoice).getId();
        log.info("Invoice created successfully with ID: {}", createdId);
        return createdId;
    }

    @Override
    public void update(UUID id, InvoiceDTO invoiceDTO) {
        log.info("Updating invoice with ID: {}", id);

        Invoice invoice = invoiceRepository.findById(id)
                .orElseThrow(() -> {
                    log.error("Invoice with ID {} not found. Update aborted.", id);
                    return new NotFoundException("Invoice not found for ID: " + id);
                });

        User owner = userRepository.findById(invoiceDTO.getOwner())
                .orElseThrow(() -> {
                    log.error("User with ID {} not found. Invoice update aborted.", invoiceDTO.getOwner());
                    return new NotFoundException("User not found for ID: " + invoiceDTO.getOwner());
                });

        Customer customer = customerRepository.findById(invoiceDTO.getCustomer())
                .orElseThrow(() -> {
                    log.error("Customer with ID {} not found. Invoice update aborted.", invoiceDTO.getCustomer());
                    return new NotFoundException("Customer not found for ID: " + invoiceDTO.getCustomer());
                });

        invoice.setAmount(invoiceDTO.getAmount());
        invoice.setDueDate(invoiceDTO.getDueDate());
        invoice.setStatus(invoiceDTO.getStatus());
        invoice.setOwner(owner);
        invoice.setCustomer(customer);
        invoiceRepository.save(invoice);

        log.info("Invoice with ID {} updated successfully.", id);
    }

    @Override
    public void delete(UUID id) {
        log.info("Deleting invoice with ID: {}", id);
        if (invoiceRepository.existsById(id)) {
            invoiceRepository.deleteById(id);
            log.info("Invoice with ID {} deleted successfully.", id);
        } else {
            log.warn("Attempted to delete non-existent invoice with ID: {}", id);
        }
    }

    @Override
    public long countByOwner(UUID ownerId) {
        log.info("Counting invoices for owner ID: {}", ownerId);
        return invoiceRepository.countByOwnerId(ownerId);
    }
}
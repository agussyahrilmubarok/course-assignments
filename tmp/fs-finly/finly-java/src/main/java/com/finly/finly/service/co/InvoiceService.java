package com.finly.finly.service.co;

import com.finly.finly.domain.Customer;
import com.finly.finly.domain.Invoice;
import com.finly.finly.domain.User;
import com.finly.finly.events.BeforeDeleteCustomer;
import com.finly.finly.events.BeforeDeleteUser;
import com.finly.finly.model.InvoiceDTO;
import com.finly.finly.repos.CustomerRepository;
import com.finly.finly.repos.InvoiceRepository;
import com.finly.finly.repos.UserRepository;
import com.finly.finly.util.NotFoundException;
import com.finly.finly.util.ReferencedException;
import org.springframework.context.event.EventListener;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;


@Service
public class InvoiceService {

    private final InvoiceRepository invoiceRepository;
    private final UserRepository userRepository;
    private final CustomerRepository customerRepository;

    public InvoiceService(final InvoiceRepository invoiceRepository,
                          final UserRepository userRepository, final CustomerRepository customerRepository) {
        this.invoiceRepository = invoiceRepository;
        this.userRepository = userRepository;
        this.customerRepository = customerRepository;
    }

    public List<InvoiceDTO> findAll() {
        final List<Invoice> invoices = invoiceRepository.findAll(Sort.by("id"));
        return invoices.stream()
                .map(invoice -> mapToDTO(invoice, new InvoiceDTO()))
                .toList();
    }

    public InvoiceDTO get(final UUID id) {
        return invoiceRepository.findById(id)
                .map(invoice -> mapToDTO(invoice, new InvoiceDTO()))
                .orElseThrow(NotFoundException::new);
    }

    public UUID create(final InvoiceDTO invoiceDTO) {
        final Invoice invoice = new Invoice();
        mapToEntity(invoiceDTO, invoice);
        return invoiceRepository.save(invoice).getId();
    }

    public void update(final UUID id, final InvoiceDTO invoiceDTO) {
        final Invoice invoice = invoiceRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        mapToEntity(invoiceDTO, invoice);
        invoiceRepository.save(invoice);
    }

    public void delete(final UUID id) {
        final Invoice invoice = invoiceRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        invoiceRepository.delete(invoice);
    }

    private InvoiceDTO mapToDTO(final Invoice invoice, final InvoiceDTO invoiceDTO) {
        invoiceDTO.setId(invoice.getId());
        invoiceDTO.setAmount(invoice.getAmount());
        invoiceDTO.setDueDate(invoice.getDueDate());
        invoiceDTO.setStatus(invoice.getStatus());
        invoiceDTO.setOwner(invoice.getOwner() == null ? null : invoice.getOwner().getId());
        invoiceDTO.setCustomer(invoice.getCustomer() == null ? null : invoice.getCustomer().getId());
        return invoiceDTO;
    }

    private Invoice mapToEntity(final InvoiceDTO invoiceDTO, final Invoice invoice) {
        invoice.setAmount(invoiceDTO.getAmount());
        invoice.setDueDate(invoiceDTO.getDueDate());
        invoice.setStatus(invoiceDTO.getStatus());
        final User owner = invoiceDTO.getOwner() == null ? null : userRepository.findById(invoiceDTO.getOwner())
                .orElseThrow(() -> new NotFoundException("owner not found"));
        invoice.setOwner(owner);
        final Customer customer = invoiceDTO.getCustomer() == null ? null : customerRepository.findById(invoiceDTO.getCustomer())
                .orElseThrow(() -> new NotFoundException("customer not found"));
        invoice.setCustomer(customer);
        return invoice;
    }

    @EventListener(BeforeDeleteUser.class)
    public void on(final BeforeDeleteUser event) {
        final ReferencedException referencedException = new ReferencedException();
        final Invoice ownerInvoice = invoiceRepository.findFirstByOwnerId(event.getId());
        if (ownerInvoice != null) {
            referencedException.setKey("user.invoice.owner.referenced");
            referencedException.addParam(ownerInvoice.getId());
            throw referencedException;
        }
    }

    @EventListener(BeforeDeleteCustomer.class)
    public void on(final BeforeDeleteCustomer event) {
        final ReferencedException referencedException = new ReferencedException();
        final Invoice customerInvoice = invoiceRepository.findFirstByCustomerId(event.getId());
        if (customerInvoice != null) {
            referencedException.setKey("customer.invoice.customer.referenced");
            referencedException.addParam(customerInvoice.getId());
            throw referencedException;
        }
    }

}

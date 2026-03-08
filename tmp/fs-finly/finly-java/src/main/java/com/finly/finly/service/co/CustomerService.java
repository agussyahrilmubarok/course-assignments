package com.finly.finly.service.co;

import com.finly.finly.domain.Customer;
import com.finly.finly.domain.User;
import com.finly.finly.events.BeforeDeleteCustomer;
import com.finly.finly.events.BeforeDeleteUser;
import com.finly.finly.model.CustomerDTO;
import com.finly.finly.repos.CustomerRepository;
import com.finly.finly.repos.UserRepository;
import com.finly.finly.util.CustomCollectors;
import com.finly.finly.util.NotFoundException;
import com.finly.finly.util.ReferencedException;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.context.event.EventListener;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Map;
import java.util.UUID;


@Service
public class CustomerService {

    private final CustomerRepository customerRepository;
    private final UserRepository userRepository;
    private final ApplicationEventPublisher publisher;

    public CustomerService(final CustomerRepository customerRepository,
                           final UserRepository userRepository, final ApplicationEventPublisher publisher) {
        this.customerRepository = customerRepository;
        this.userRepository = userRepository;
        this.publisher = publisher;
    }

    public List<CustomerDTO> findAll() {
        final List<Customer> customers = customerRepository.findAll(Sort.by("id"));
        return customers.stream()
                .map(customer -> mapToDTO(customer, new CustomerDTO()))
                .toList();
    }

    public CustomerDTO get(final UUID id) {
        return customerRepository.findById(id)
                .map(customer -> mapToDTO(customer, new CustomerDTO()))
                .orElseThrow(NotFoundException::new);
    }

    public UUID create(final CustomerDTO customerDTO) {
        final Customer customer = new Customer();
        mapToEntity(customerDTO, customer);
        return customerRepository.save(customer).getId();
    }

    public void update(final UUID id, final CustomerDTO customerDTO) {
        final Customer customer = customerRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        mapToEntity(customerDTO, customer);
        customerRepository.save(customer);
    }

    public void delete(final UUID id) {
        final Customer customer = customerRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        publisher.publishEvent(new BeforeDeleteCustomer(id));
        customerRepository.delete(customer);
    }

    private CustomerDTO mapToDTO(final Customer customer, final CustomerDTO customerDTO) {
        customerDTO.setId(customer.getId());
        customerDTO.setName(customer.getName());
        customerDTO.setEmail(customer.getEmail());
        customerDTO.setPhone(customer.getPhone());
        customerDTO.setAddress(customer.getAddress());
        customerDTO.setUser(customer.getUser() == null ? null : customer.getUser().getId());
        return customerDTO;
    }

    private Customer mapToEntity(final CustomerDTO customerDTO, final Customer customer) {
        customer.setName(customerDTO.getName());
        customer.setEmail(customerDTO.getEmail());
        customer.setPhone(customerDTO.getPhone());
        customer.setAddress(customerDTO.getAddress());
        final User user = customerDTO.getUser() == null ? null : userRepository.findById(customerDTO.getUser())
                .orElseThrow(() -> new NotFoundException("user not found"));
        customer.setUser(user);
        return customer;
    }

    public boolean emailExists(final String email) {
        return customerRepository.existsByEmailIgnoreCase(email);
    }

    public Map<UUID, String> getCustomerValues() {
        return customerRepository.findAll(Sort.by("id"))
                .stream()
                .collect(CustomCollectors.toSortedMap(Customer::getId, Customer::getName));
    }

    @EventListener(BeforeDeleteUser.class)
    public void on(final BeforeDeleteUser event) {
        final ReferencedException referencedException = new ReferencedException();
        final Customer userCustomer = customerRepository.findFirstByUserId(event.getId());
        if (userCustomer != null) {
            referencedException.setKey("user.customer.user.referenced");
            referencedException.addParam(userCustomer.getId());
            throw referencedException;
        }
    }

}

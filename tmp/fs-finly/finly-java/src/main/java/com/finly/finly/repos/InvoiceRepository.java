package com.finly.finly.repos;

import com.finly.finly.domain.Invoice;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.data.mongodb.repository.Query;

import java.util.List;
import java.util.UUID;


public interface InvoiceRepository extends MongoRepository<Invoice, UUID> {

    @Query("{ 'owner': ?0, 'customer.name': { $regex: ?1, $options: 'i' } }")
    List<Invoice> searchInvoicesByOwnerId(UUID ownerId, String search);

    List<Invoice> findByOwnerId(UUID ownerId);

    Invoice findFirstByOwnerId(UUID id);

    Invoice findFirstByCustomerId(UUID id);

    long countByOwnerId(UUID ownerId);
}

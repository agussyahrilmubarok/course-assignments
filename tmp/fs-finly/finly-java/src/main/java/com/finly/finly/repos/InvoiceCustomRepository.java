package com.finly.finly.repos;

import com.finly.finly.domain.Invoice;

import java.util.List;
import java.util.UUID;

public interface InvoiceCustomRepository {

    List<Invoice> findInvoicesByOwnerAndSearch(UUID ownerId, String search);

    List<Invoice> findInvoicesByOwnerId(UUID ownerId);
}

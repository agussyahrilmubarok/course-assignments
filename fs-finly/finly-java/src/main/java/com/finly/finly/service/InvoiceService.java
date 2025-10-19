package com.finly.finly.service;

import com.finly.finly.model.InvoiceDTO;
import com.finly.finly.model.InvoiceDetailDTO;

import java.util.List;
import java.util.UUID;

public interface InvoiceService {

    List<InvoiceDTO> findAll();

    List<InvoiceDetailDTO> findInvoiceByOwnerAndSearch(UUID userId, String search);

    InvoiceDTO get(UUID id);

    InvoiceDetailDTO getDetail(UUID id);

    UUID create(InvoiceDTO invoiceDTO);

    void update(UUID id, InvoiceDTO invoiceDTO);

    void delete(UUID id);

    long countByOwner(UUID ownerId);
}

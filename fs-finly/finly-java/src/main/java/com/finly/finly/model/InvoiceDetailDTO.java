package com.finly.finly.model;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.finly.finly.domain.Invoice;
import jakarta.validation.constraints.Digits;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.util.UUID;

@Getter
@Setter
public class InvoiceDetailDTO {

    private UUID id;

    @NotNull
    @Digits(integer = 10, fraction = 2)
    @JsonFormat(shape = JsonFormat.Shape.STRING)
    private BigDecimal amount;

    @NotNull
    @JsonFormat(pattern = "MM/dd/yyyy")
    private LocalDate dueDate;

    @NotNull
    @Size(max = 255)
    private String status;

    @NotNull
    private CustomerDTO customer;

    private UserDTO owner;

    public static InvoiceDetailDTO fromInvoice(Invoice invoice) {
        InvoiceDetailDTO invoiceDTO = new InvoiceDetailDTO();
        invoiceDTO.setId(invoice.getId());
        invoiceDTO.setAmount(invoice.getAmount());
        invoiceDTO.setDueDate(invoice.getDueDate());
        invoiceDTO.setStatus(invoice.getStatus());
        invoiceDTO.setCustomer(CustomerDTO.fromCustomer(invoice.getCustomer()));
        invoiceDTO.setOwner(UserDTO.fromUser(invoice.getOwner()));
        return invoiceDTO;
    }
}

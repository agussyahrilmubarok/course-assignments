package com.finly.finly.model;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.finly.finly.domain.Invoice;
import jakarta.validation.constraints.Digits;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;
import org.springframework.format.annotation.DateTimeFormat;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.util.UUID;


@Getter
@Setter
public class InvoiceDTO {

    private UUID id;

    @NotNull
    @Digits(integer = 10, fraction = 2)
    @JsonFormat(shape = JsonFormat.Shape.STRING)
    private BigDecimal amount;

    @NotNull
    @DateTimeFormat(pattern = "MM/dd/yyyy")
    @JsonFormat(pattern = "MM/dd/yyyy")
    private LocalDate dueDate;

    @NotNull
    @Size(max = 255)
    private String status;

    @NotNull
    private UUID customer;

    private UUID owner;

    public static InvoiceDTO fromInvoice(Invoice invoice) {
        InvoiceDTO invoiceDTO = new InvoiceDTO();
        invoiceDTO.setId(invoice.getId());
        invoiceDTO.setAmount(invoice.getAmount());
        invoiceDTO.setDueDate(invoice.getDueDate());
        invoiceDTO.setStatus(invoice.getStatus());
        invoiceDTO.setCustomer(invoice.getCustomer().getId());
        invoiceDTO.setOwner(invoice.getOwner().getId());
        return invoiceDTO;
    }
}

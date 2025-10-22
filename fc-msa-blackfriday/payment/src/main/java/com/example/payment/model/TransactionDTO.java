package com.example.payment.model;

import com.example.payment.domain.Transaction;
import com.fasterxml.jackson.annotation.JsonInclude;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.*;

import java.math.BigDecimal;

@Getter
@Setter
public class TransactionDTO {

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class CreateTransactionRequest {

        @NotBlank(message = "Order ID must not be blank")
        private String orderId;

        @NotNull(message = "Amount is required")
        @Min(value = 0, message = "Amount must be >= 0")
        private BigDecimal amount;
    }

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class Response {
        private String id;
        private String orderId;
        private BigDecimal amount;
        private String paymentUrl;
        private String status;

        public static Response from(Transaction transaction) {
            return Response.builder()
                    .id(transaction.getId())
                    .orderId(transaction.getOrderId())
                    .amount(transaction.getAmount())
                    .paymentUrl(transaction.getPaymentUrl())
                    .status(transaction.getStatus().name())
                    .build();
        }
    }
}

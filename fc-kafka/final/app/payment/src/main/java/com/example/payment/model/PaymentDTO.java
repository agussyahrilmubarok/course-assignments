package com.example.payment.model;

import com.fasterxml.jackson.annotation.JsonFormat;
import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.Digits;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;
import org.springframework.format.annotation.DateTimeFormat;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;


@Getter
@Setter
public class PaymentDTO {

    private UUID id;

    @NotNull
    @Size(max = 255)
    private String orderId;

    @NotNull
    @Size(max = 255)
    private String customerId;

    @NotNull
    @Digits(integer = 10, fraction = 2)
    @JsonFormat(shape = JsonFormat.Shape.STRING)
    @Schema(type = "string", example = "92.08")
    private BigDecimal amount;

    @Size(max = 255)
    private String method;

    @NotNull
    @Size(max = 255)
    private String status;

    @DateTimeFormat(pattern = "yyyy-MM-dd'T'HH:mm")
    private LocalDateTime paidAt;

}

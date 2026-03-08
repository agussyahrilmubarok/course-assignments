package com.finly.finly.domain;

import jakarta.validation.constraints.Digits;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.Id;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.annotation.Version;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.DocumentReference;
import org.springframework.data.mongodb.core.mapping.Field;
import org.springframework.data.mongodb.core.mapping.FieldType;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.UUID;


@Document("invoices")
@Getter
@Setter
public class Invoice {

    @Id
    private UUID id;

    @NotNull
    @Digits(integer = 10, fraction = 2)
    @Field(
            targetType = FieldType.DECIMAL128
    )
    private BigDecimal amount;

    @NotNull
    private LocalDate dueDate;

    @NotNull
    @Size(max = 255)
    private String status;

    @DocumentReference(lazy = true)
    @NotNull
    private User owner;

    @DocumentReference(lazy = true)
    @NotNull
    private Customer customer;

    @CreatedDate
    private OffsetDateTime createdAt;

    @LastModifiedDate
    private OffsetDateTime updatedAt;

    @Version
    private Integer version;
}

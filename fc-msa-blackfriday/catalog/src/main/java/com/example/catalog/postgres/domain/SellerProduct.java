package com.example.catalog.postgres.domain;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "SellerProducts")
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
public class SellerProduct {

    @Id
    @Column(nullable = false, updatable = false)
    private String id;

    @Column(nullable = false)
    public String sellerId;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private LocalDateTime updatedAt;
}

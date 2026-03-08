package com.example.point.domain;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "PointBalances")
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
public class PointBalance {

    @Id
    @Column(nullable = false, updatable = false)
    private String id;

    @Column(nullable = false)
    private Long balance = 0L;

    @Column(nullable = false)
    private String userId;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private LocalDateTime updatedAt;

    @Version
    private Long version = 0L;

    public void addBalance(Long amount) {
        if (amount <= 0) {
            throw new IllegalArgumentException("Amount must be positive");
        }
        if (this.balance == null) {
            this.balance = 0L;
        }
        this.balance += amount;
    }

    public void subtractBalance(Long amount) {
        if (amount <= 0) {
            throw new IllegalArgumentException("Amount must be positive");
        }
        if (this.balance == null) {
            this.balance = 0L;
        }
        if (this.balance < amount) {
            throw new IllegalArgumentException("Insufficient point balance");
        }
        this.balance -= amount;
    }

    public void setBalance(Long balance) {
        if (balance == null || balance < 0) {
            throw new IllegalArgumentException("Balance cannot be negative or null");
        }
        this.balance = balance;
    }
}

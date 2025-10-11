package com.example.point.batch.point.domain;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "Points")
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
public class Point {

    @Id
    @Column(nullable = false, updatable = false)
    private String id;

    @Column(nullable = false)
    private Long amount;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private PointType type;

    @Column(nullable = false)
    private String description;

    @Column(nullable = false)
    private Long balanceSnapshot;

    @Column(nullable = false)
    private String userId;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "point_balance_id", nullable = false)
    private PointBalance pointBalance;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private LocalDateTime updatedAt;

    @Version
    private Long version = 0L;

    public enum PointType {
        EARNED,
        USED,
        CANCELED
    }
}

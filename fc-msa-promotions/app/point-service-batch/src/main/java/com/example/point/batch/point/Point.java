package com.example.point.batch.point;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.time.LocalDateTime;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class Point implements Serializable {

    private String id;
    private Long amount;
    private PointType type;
    private String description;
    private Long balanceSnapshot;
    private String userId;
    private String pointBalanceId;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
    private Long version = 0L;

    public enum PointType {
        EARNED,
        USED,
        CANCELED
    }
}

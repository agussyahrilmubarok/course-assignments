package com.example.point.batch.point;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.time.LocalDateTime;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class PointBalance implements Serializable {

    private String id;
    private Long balance = 0L;
    private String userId;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
    private Long version = 0L;
}

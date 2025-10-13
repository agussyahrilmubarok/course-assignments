package com.example.point.scheduler;

import com.example.point.domain.Point;
import com.example.point.domain.Point.PointType;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class PointRedisDTO {

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

    public static PointRedisDTO toDto(Point point) {
        return PointRedisDTO.builder()
                .id(point.getId())
                .amount(point.getAmount())
                .type(point.getType())
                .description(point.getDescription())
                .balanceSnapshot(point.getBalanceSnapshot())
                .userId(point.getUserId())
                .pointBalanceId(point.getPointBalance().getId())
                .createdAt(point.getCreatedAt())
                .updatedAt(point.getUpdatedAt())
                .version(point.getVersion())
                .build();
    }
}

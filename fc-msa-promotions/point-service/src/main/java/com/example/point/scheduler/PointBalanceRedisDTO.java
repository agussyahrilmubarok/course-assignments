package com.example.point.scheduler;

import com.example.point.domain.PointBalance;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class PointBalanceRedisDTO {

    private String id;
    private Long balance = 0L;
    private String userId;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
    private Long version = 0L;

    public static PointBalanceRedisDTO toDto(PointBalance balance) {
        return PointBalanceRedisDTO.builder()
                .id(balance.getId())
                .balance(balance.getBalance())
                .userId(balance.getUserId())
                .createdAt(balance.getCreatedAt())
                .updatedAt(balance.getUpdatedAt())
                .version(balance.getVersion())
                .build();
    }
}


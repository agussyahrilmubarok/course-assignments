package com.example.point.model;

import com.example.point.domain.Point;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

import java.time.LocalDateTime;

@Getter
@Setter
public class PointDTO {

    @Getter
    @Builder
    public static class EarnRequest {
        private String userId;

        @NotNull(message = "Amount must not be null")
        @Min(value = 1, message = "Amount must be greater than 0")
        private Long amount;

        @NotBlank(message = "Description must not be blank")
        private String description;
    }

    @Getter
    @Builder
    public static class UseRequest {
        private String userId;

        @NotNull(message = "Amount must not be null")
        @Min(value = 1, message = "Amount must be greater than 0")
        private Long amount;

        @NotBlank(message = "Description must not be blank")
        private String description;
    }

    @Getter
    @Builder
    public static class CancelRequest {
        @NotNull(message = "Point id must not be null")
        private String pointId;

        @NotBlank(message = "Description must not be blank")
        private String description;
    }

    @Getter
    @Builder
    public static class Response {
        private String id;
        private String userId;
        private Long amount;
        private Point.PointType type;
        private String description;
        private Long balanceSnapshot;
        private LocalDateTime createdAt;

        public static Response from(Point point) {
            return Response.builder()
                    .id(point.getId())
                    .userId(point.getUserId())
                    .amount(point.getAmount())
                    .type(point.getType())
                    .description(point.getDescription())
                    .balanceSnapshot(point.getBalanceSnapshot())
                    .createdAt(point.getCreatedAt())
                    .build();
        }
    }

    @Getter
    @Builder
    public static class BalanceResponse {
        private String userId;
        private Long balance;

        public static BalanceResponse of(String userId, Long balance) {
            return BalanceResponse.builder()
                    .userId(userId)
                    .balance(balance)
                    .build();
        }
    }
}

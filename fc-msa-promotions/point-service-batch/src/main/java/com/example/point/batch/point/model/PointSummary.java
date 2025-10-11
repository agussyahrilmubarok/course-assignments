package com.example.point.batch.point.model;

import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Getter
@NoArgsConstructor
public class PointSummary {

    private String userId;
    private Long earnAmount;
    private Long useAmount;
    private Long cancelAmount;

    @Builder
    public PointSummary(String userId, Long earnAmount, Long useAmount, Long cancelAmount) {
        this.userId = userId;
        this.earnAmount = earnAmount;
        this.useAmount = useAmount;
        this.cancelAmount = cancelAmount;
    }

    public void addEarnAmount(Long amount) {
        this.earnAmount += amount;
    }

    public void addUseAmount(Long amount) {
        this.useAmount += amount;
    }

    public void addCancelAmount(Long amount) {
        this.cancelAmount += amount;
    }
}

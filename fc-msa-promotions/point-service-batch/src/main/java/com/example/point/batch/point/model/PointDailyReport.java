package com.example.point.batch.point.model;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

import java.time.LocalDate;

@Getter
@Setter
public class PointDailyReport {

    private String id;
    private String userId;
    private LocalDate reportDate;
    private Long earnAmount;
    private Long useAmount;
    private Long cancelAmount;
    private Long netAmount;

    @Builder
    public PointDailyReport(String userId, LocalDate reportDate, Long earnAmount,
                            Long useAmount, Long cancelAmount) {
        this.userId = userId;
        this.reportDate = reportDate;
        this.earnAmount = earnAmount;
        this.useAmount = useAmount;
        this.cancelAmount = cancelAmount;
        this.netAmount = earnAmount - useAmount + cancelAmount;
    }

}

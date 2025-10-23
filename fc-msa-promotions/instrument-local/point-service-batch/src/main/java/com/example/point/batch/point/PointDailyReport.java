package com.example.point.batch.point;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDate;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class PointDailyReport {

    private Long userId;
    private LocalDate reportDate;
    private Long earnAmount;
    private Long useAmount;
    private Long cancelAmount;
    private Long netAmount;

    public void calculateNetAmount() {
        this.netAmount = (earnAmount != null ? earnAmount : 0L)
                - (useAmount != null ? useAmount : 0L)
                + (cancelAmount != null ? cancelAmount : 0L);
    }
}

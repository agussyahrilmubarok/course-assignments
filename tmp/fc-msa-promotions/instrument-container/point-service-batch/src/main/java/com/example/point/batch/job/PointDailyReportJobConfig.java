package com.example.point.batch.job;

import com.example.point.batch.listener.JobCompletionNotificationListener;
import com.example.point.batch.point.Point;
import com.example.point.batch.point.PointBalance;
import com.example.point.batch.point.PointSummary;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RBucket;
import org.redisson.api.RedissonClient;
import org.springframework.batch.core.Job;
import org.springframework.batch.core.Step;
import org.springframework.batch.core.configuration.annotation.StepScope;
import org.springframework.batch.core.job.builder.JobBuilder;
import org.springframework.batch.core.repository.JobRepository;
import org.springframework.batch.core.step.builder.StepBuilder;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.batch.item.ItemReader;
import org.springframework.batch.item.ItemWriter;
import org.springframework.batch.item.support.IteratorItemReader;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.transaction.PlatformTransactionManager;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.*;
import java.util.Map.Entry;

import static java.time.LocalDate.now;

@Configuration
@Slf4j
@RequiredArgsConstructor
public class PointDailyReportJobConfig {

    private static final String BALANCE_PREFIX = "cache:balance:";
    private static final String POINT_PREFIX = "cache:point:";
    private final JobRepository jobRepository;
    private final PlatformTransactionManager transactionManager;
    private final JobCompletionNotificationListener jobListener;
    private final RedissonClient redissonClient;
    private final ObjectMapper objectMapper;

    @Bean
    public Job pointDailyReportJob() {
        return new JobBuilder("pointDailyReportJob", jobRepository)
                .listener(jobListener)
                .start(syncPointBalanceStep())
                .next(generateDailyReportStep())
                .build();
    }

    @Bean
    public Step syncPointBalanceStep() {
        return new StepBuilder("syncPointBalanceStep", jobRepository)
                .<PointBalance, Entry<String, Long>>chunk(1000, transactionManager)
                .reader(pointBalanceReader())
                .processor(pointBalanceProcessor())
                .writer(pointBalanceWriter())
                .build();
    }

    @Bean
    @StepScope
    public ItemReader<PointBalance> pointBalanceReader() {
        Iterable<String> keys = redissonClient.getKeys().getKeysByPattern(BALANCE_PREFIX + "*");
        List<PointBalance> balances = new ArrayList<>();
        for (String key : keys) {
            RBucket<String> bucket = redissonClient.getBucket(key);
            String json = bucket.get();
            if (json != null) {
                try {
                    PointBalance balance = objectMapper.readValue(json, PointBalance.class);
                    balances.add(balance);
                } catch (Exception e) {
                    log.error("Failed to parse PointBalance for key {}: {}", key, e.getMessage());
                }
            }
        }
        return new IteratorItemReader<>(balances);
    }

    @Bean
    @StepScope
    public ItemProcessor<PointBalance, Entry<String, Long>> pointBalanceProcessor() {
        return balance -> Map.entry(balance.getUserId(), balance.getBalance());
    }

    @Bean
    @StepScope
    public ItemWriter<Entry<String, Long>> pointBalanceWriter() {
        return items -> {
            for (Entry<String, Long> entry : items) {
                log.info("Synced Balance - userId: {}, balance: {}", entry.getKey(), entry.getValue());
            }
        };
    }

    @Bean
    public Step generateDailyReportStep() {
        return new StepBuilder("generateDailyReportStep", jobRepository)
                .<Point, PointSummary>chunk(1000, transactionManager)
                .reader(pointReader())
                .processor(pointProcessor())
                .writer(reportWriter())
                .build();
    }

    @Bean
    @StepScope
    public ItemReader<Point> pointReader() {
        return new ItemReader<>() {
            private Iterator<Point> iterator;

            @Override
            public Point read() {
                if (iterator == null) {
                    Iterable<String> keys = redissonClient.getKeys().getKeysByPattern(POINT_PREFIX + "*");

                    if (keys == null) return null;
                    // If run on 1 AM set to yesterday
                    LocalDateTime now = LocalDateTime.now();
                    LocalDateTime startTime = now.withHour(0).withMinute(0).withSecond(0).withNano(0);
                    LocalDateTime endTime = now.withHour(23).withMinute(59).withSecond(59).withNano(999_999_999);

                    List<Point> filtered = new ArrayList<>();

                    for (String key : keys) {
                        RBucket<String> bucket = redissonClient.getBucket(key);
                        String json = bucket.get();
                        if (json != null) {
                            try {
                                Point point = objectMapper.readValue(json, Point.class);
                                if (point.getCreatedAt() != null &&
                                        !point.getCreatedAt().isBefore(startTime) &&
                                        !point.getCreatedAt().isAfter(endTime)) {
                                    filtered.add(point);
                                }
                            } catch (Exception e) {
                                log.error("Failed to parse Point for key {}: {}", key, e.getMessage());
                            }
                        }
                    }

                    iterator = filtered.iterator();
                }

                return iterator.hasNext() ? iterator.next() : null;
            }
        };
    }

    @Bean
    @StepScope
    public ItemProcessor<Point, PointSummary> pointProcessor() {
        return point -> {
            switch (point.getType()) {
                case EARNED -> {
                    return new PointSummary(point.getUserId(), point.getAmount(), 0L, 0L);
                }
                case USED -> {
                    return new PointSummary(point.getUserId(), 0L, point.getAmount(), 0L);
                }
                case CANCELED -> {
                    return new PointSummary(point.getUserId(), 0L, 0L, point.getAmount());
                }
                default -> {
                    return null;
                }
            }
        };
    }

    @Bean
    @StepScope
    public ItemWriter<PointSummary> reportWriter() {
        return items -> {
            String exportDir = "exports";
            String fileName = "daily-report-" + now().format(DateTimeFormatter.ISO_DATE) + ".csv";
            String fullPath = exportDir + "/" + fileName;

            File dir = new File(exportDir);
            if (!dir.exists() && !dir.mkdirs()) {
                throw new IOException("Failed to create export directory: " + exportDir);
            }

            // Aggregate PointSummary by userId
            Map<String, PointSummary> summaryMap = new HashMap<>();
            for (PointSummary item : items) {
                String userId = item.getUserId();
                PointSummary existing = summaryMap.getOrDefault(userId, new PointSummary(userId, 0L, 0L, 0L));

                existing.setEarnAmount(existing.getEarnAmount() + (item.getEarnAmount() != null ? item.getEarnAmount() : 0L));
                existing.setUseAmount(existing.getUseAmount() + (item.getUseAmount() != null ? item.getUseAmount() : 0L));
                existing.setCancelAmount(existing.getCancelAmount() + (item.getCancelAmount() != null ? item.getCancelAmount() : 0L));

                summaryMap.put(userId, existing);
            }

            try (FileWriter writer = new FileWriter(fullPath)) {
                writer.write("userId,reportDate,earnAmount,useAmount,cancelAmount,netAmount\n");
                for (PointSummary summary : summaryMap.values()) {
                    long earn = summary.getEarnAmount() != null ? summary.getEarnAmount() : 0L;
                    long use = summary.getUseAmount() != null ? summary.getUseAmount() : 0L;
                    long cancel = summary.getCancelAmount() != null ? summary.getCancelAmount() : 0L;
                    long net = earn - use + cancel;

                    writer.write(String.format(
                            "%s,%s,%d,%d,%d,%d\n",
                            summary.getUserId(),
                            now().minusDays(1), // tanggal laporan
                            earn,
                            use,
                            cancel,
                            net
                    ));
                }

                log.info("Exported {} rows to {}", summaryMap.size(), fullPath);
            } catch (IOException e) {
                log.error("Failed to write CSV file to {}: {}", fullPath, e.getMessage(), e);
                throw new RuntimeException(e);
            }
        };
    }
}

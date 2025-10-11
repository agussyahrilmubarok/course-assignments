package com.example.point.batch.job;

import com.example.point.batch.point.domain.Point;
import com.example.point.batch.point.domain.PointBalance;
import com.example.point.batch.listener.JobCompletionNotificationListener;
import com.example.point.batch.point.model.PointDailyReport;
import com.example.point.batch.point.model.PointSummary;
import jakarta.persistence.EntityManagerFactory;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RedissonClient;
import org.springframework.batch.core.Job;
import org.springframework.batch.core.Step;
import org.springframework.batch.core.configuration.annotation.StepScope;
import org.springframework.batch.core.job.builder.JobBuilder;
import org.springframework.batch.core.repository.JobRepository;
import org.springframework.batch.core.step.builder.StepBuilder;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.batch.item.ItemWriter;
import org.springframework.batch.item.database.JpaPagingItemReader;
import org.springframework.batch.item.database.builder.JpaPagingItemReaderBuilder;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.transaction.PlatformTransactionManager;

import java.io.FileWriter;
import java.io.PrintWriter;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static java.time.LocalDate.now;

/**
 * Batch Job configuration for synchronizing point balances and generating daily reports
 * <p>
 * Main Features:
 * 1. Synchronize point balances between Redis cache and database
 * 2. Generate daily reports based on the previous day's point transactions
 */
@Configuration
@Slf4j
@RequiredArgsConstructor
public class PointBalanceSyncJobConfig {

    private final JobRepository jobRepository;
    @Qualifier("pointTransactionManager")
    private final PlatformTransactionManager transactionManager;
    private final JobCompletionNotificationListener jobCompletionNotificationListener;
    @Qualifier("pointEntityManagerFactory")
    private final EntityManagerFactory entityManagerFactory;
    private final RedissonClient redissonClient;

    /**
     * Job to synchronize point balances and generate daily reports
     * <p>
     * Execution order:
     * 1. syncPointBalanceStep: Synchronize point balances from DB to Redis cache
     * 2. generateDailyReportStep: Aggregate previous day's point transactions and create daily reports
     */
    @Bean
    public Job pointBalanceSyncJob() {
        return new JobBuilder("pointBalanceSyncJob", jobRepository)
                .listener(jobCompletionNotificationListener)
                .start(syncPointBalanceStep())
                .next(generateDailyReportStep())
                .build();
    }

    /**
     * Step to synchronize point balances
     * <p>
     * Synchronizes point balance data from the database to Redis cache
     * - Reader: Read point balances using JPA
     * - Processor: Generate cache keys
     * - Writer: Save point balances to Redis
     */
    @Bean
    public Step syncPointBalanceStep() {
        return new StepBuilder("syncPointBalanceStep", jobRepository)
                .<PointBalance, Map.Entry<String, Long>>chunk(1000, transactionManager)
                .reader(pointBalanceReader())
                .processor(pointBalanceProcessor())
                .writer(pointBalanceWriter())
                .build();
    }

    /**
     * Step to generate daily report
     * <p>
     * Aggregates point transactions from the previous day to create daily reports
     * - Reader: Read previous day's point transactions using JPA
     * - Processor: Aggregate transactions per user
     * - Writer: Save daily reports to the database
     */
    @Bean
    public Step generateDailyReportStep() {
        return new StepBuilder("generateDailyReportStep", jobRepository)
                .<Point, PointSummary>chunk(1000, transactionManager)
                .reader(pointReader())
                .processor(pointProcessor())
                .writer(csvReportWriter())
                .build();
    }

    /**
     * Reader for point balances
     * <p>
     * Reads point balance information from the database using JPA
     */
    @Bean
    @StepScope
    public JpaPagingItemReader<PointBalance> pointBalanceReader() {
        return new JpaPagingItemReaderBuilder<PointBalance>()
                .name("pointBalanceReader")
                .entityManagerFactory(entityManagerFactory)
                .pageSize(1000)
                .queryString("SELECT pb FROM PointBalance pb")
                .build();
    }

    /**
     * Processor for point balances
     * <p>
     * Converts point balances into key-value pairs for Redis caching
     */
    @Bean
    @StepScope
    public ItemProcessor<PointBalance, Map.Entry<String, Long>> pointBalanceProcessor() {
        return pointBalance -> Map.entry(
                String.format("point:balance:%s", pointBalance.getUserId()),
                pointBalance.getBalance()
        );
    }

    /**
     * Writer for point balances
     * <p>
     * Saves point balances to Redis cache
     */
    @Bean
    @StepScope
    public ItemWriter<Map.Entry<String, Long>> pointBalanceWriter() {
        return items -> {
            var balanceMap = redissonClient.getMap("point:balance");
            items.forEach(item -> balanceMap.put(item.getKey(), item.getValue()));
        };
    }

    /**
     * Reader for point transactions
     * <p>
     * Reads point transactions from the previous day
     */
    @Bean
    @StepScope
    public JpaPagingItemReader<Point> pointReader() {
        Map<String, Object> parameters = new HashMap<>();
        LocalDateTime yesterday = LocalDateTime.now().minusDays(1);
        parameters.put("startTime", yesterday.withHour(0).withMinute(0).withSecond(0));
        parameters.put("endTime", yesterday.withHour(23).withMinute(59).withSecond(59));

        return new JpaPagingItemReaderBuilder<Point>()
                .name("pointReader")
                .entityManagerFactory(entityManagerFactory)
                .pageSize(1000)
                .queryString("SELECT p FROM Point p WHERE p.createdAt BETWEEN :startTime AND :endTime")
                .parameterValues(parameters)
                .build();
    }

    /**
     * Processor for point transactions
     * <p>
     * Aggregates point transactions per user and generates PointSummary objects
     */
    @Bean
    @StepScope
    public ItemProcessor<Point, PointSummary> pointProcessor() {
        return point -> {
            return switch (point.getType()) {
                case EARNED -> new PointSummary(point.getUserId(), point.getAmount(), 0L, 0L);
                case USED -> new PointSummary(point.getUserId(), 0L, point.getAmount(), 0L);
                case CANCELED -> new PointSummary(point.getUserId(), 0L, 0L, point.getAmount());
            };
        };
    }

    /**
     * Writer for daily reports
     * <p>
     * Converts aggregated point transactions into daily reports and saves them to the database
     */
    @Bean
    @StepScope
    public ItemWriter<PointSummary> dbReportWriter() {
        return summaries -> {
            List<PointDailyReport> reports = new ArrayList<>();
            for (PointSummary summary : summaries) {
                PointDailyReport report = PointDailyReport.builder()
                        .userId(summary.getUserId())
                        .reportDate(now().minusDays(1))  // Previous day's data
                        .earnAmount(summary.getEarnAmount())
                        .useAmount(summary.getUseAmount())
                        .cancelAmount(summary.getCancelAmount())
                        .build();
                reports.add(report);
            }
            // Write to storage
            // dailyPointReportRepository.saveAll(reports);
        };
    }

    @Bean
    @StepScope
    public ItemWriter<PointSummary> csvReportWriter() {
        return summaries -> {
            String fileName = String.format("exports/daily_report_%s.csv", now().minusDays(1));
            try (PrintWriter writer = new PrintWriter(new FileWriter(fileName))) {
                writer.println("userId,reportDate,earnAmount,useAmount,cancelAmount");
                for (PointSummary summary : summaries) {
                    writer.printf(
                            "%s,%s,%d,%d,%d%n",
                            summary.getUserId(),
                            now().minusDays(1),
                            summary.getEarnAmount(),
                            summary.getUseAmount(),
                            summary.getCancelAmount()
                    );
                }
                log.info("Daily report CSV exported to {}", fileName);
            } catch (Exception e) {
                log.error("Failed to write CSV report", e);
                throw new RuntimeException("CSV export failed", e);
            }
        };
    }
}


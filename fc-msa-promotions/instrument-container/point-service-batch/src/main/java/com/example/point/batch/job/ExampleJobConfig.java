package com.example.point.batch.job;

import com.example.point.batch.listener.JobCompletionNotificationListener;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.batch.core.Job;
import org.springframework.batch.core.Step;
import org.springframework.batch.core.configuration.annotation.StepScope;
import org.springframework.batch.core.job.builder.JobBuilder;
import org.springframework.batch.core.repository.JobRepository;
import org.springframework.batch.core.step.builder.StepBuilder;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.batch.item.ItemReader;
import org.springframework.batch.item.ItemWriter;
import org.springframework.batch.item.support.ListItemReader;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.transaction.PlatformTransactionManager;

import java.util.List;

/**
 * ExampleJobConfig
 * <p>
 * A simple example of a Spring Batch job that uses only an in-memory H2 database
 * for the Batch JobRepository. This job demonstrates a basic chunk-oriented
 * step pipeline:
 * 1. Read a list of names from memory
 * 2. Process each name (convert to uppercase)
 * 3. Write the results to the log
 */
@Configuration
@Slf4j
@RequiredArgsConstructor
public class ExampleJobConfig {

    private final JobRepository jobRepository;
    private final PlatformTransactionManager transactionManager;
    private final JobCompletionNotificationListener jobListener;

    /**
     * Main batch job definition
     */
    @Bean
    public Job exampleJob() {
        return new JobBuilder("exampleJob", jobRepository)
                .listener(jobListener)
                .start(exampleStep())
                .build();
    }

    /**
     * Step definition:
     * Reads a list of names → processes them → logs the output
     */
    @Bean
    public Step exampleStep() {
        return new StepBuilder("exampleStep", jobRepository)
                .<String, String>chunk(5, transactionManager)
                .reader(exampleReader())
                .processor(exampleProcessor())
                .writer(exampleWriter())
                .build();
    }

    /**
     * Reader — provides a list of names from memory
     */
    @Bean
    @StepScope
    public ItemReader<String> exampleReader() {
        List<String> data = List.of("a", "Bb", "Ccc", "Dddd", "e");
        return new ListItemReader<>(data);
    }

    /**
     * Processor — converts each name to uppercase
     */
    @Bean
    @StepScope
    public ItemProcessor<String, String> exampleProcessor() {
        return item -> item.toUpperCase();
    }

    /**
     * Writer — logs each processed item
     */
    @Bean
    @StepScope
    public ItemWriter<String> exampleWriter() {
        return items -> items.forEach(item -> log.info("Processed item: {}", item));
    }
}

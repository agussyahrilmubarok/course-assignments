package com.example.point.batch;

import org.springframework.batch.core.Job;
import org.springframework.batch.core.JobParametersBuilder;
import org.springframework.batch.core.configuration.annotation.EnableBatchProcessing;
import org.springframework.batch.core.launch.JobLauncher;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

import java.util.Map;

@SpringBootApplication
@EnableBatchProcessing
public class PointServiceBatchApplication {

    private final JobLauncher jobLauncher;
    private final Map<String, Job> jobs;

    public PointServiceBatchApplication(JobLauncher jobLauncher, Map<String, Job> jobs) {
        this.jobLauncher = jobLauncher;
        this.jobs = jobs;
    }

    public static void main(String[] args) {
        SpringApplication.run(PointServiceBatchApplication.class, args);
    }

    @Bean
    public ApplicationRunner runner() {
        return args -> {
            for (Job job : jobs.values()) {
                jobLauncher.run(
                        job,
                        new JobParametersBuilder()
                                .addLong("timestamp", System.currentTimeMillis())
                                .toJobParameters()
                );
            }
            ;
        };
    }
}

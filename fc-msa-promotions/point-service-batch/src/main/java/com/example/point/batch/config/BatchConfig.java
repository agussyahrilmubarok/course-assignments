package com.example.point.batch.config;

import org.springframework.batch.core.configuration.annotation.EnableBatchProcessing;
import org.springframework.batch.core.configuration.support.DefaultBatchConfiguration;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.jdbc.DataSourceBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.jdbc.datasource.DataSourceTransactionManager;
import org.springframework.transaction.PlatformTransactionManager;
import org.springframework.transaction.annotation.EnableTransactionManagement;

import javax.sql.DataSource;

@Configuration
@EnableBatchProcessing
public class BatchConfig extends DefaultBatchConfiguration {

    @Bean
    @ConfigurationProperties("spring.datasource.batch")
    public DataSource batchDataSource() {
        return DataSourceBuilder.create().build();
    }

    @Bean
    @ConfigurationProperties("spring.datasource.batch")
    public DataSource batchDataSource() {
        return DataSourceBuilder.create().build();
    }

    @Bean
    public PlatformTransactionManager transactionManager() {
        return new DataSourceTransactionManager(batchDataSource());
    }

    @Override
    protected DataSource getDataSource() {
        return batchDataSource();
    }

    @Override
    protected PlatformTransactionManager getTransactionManager() {
        return transactionManager();
    }
}

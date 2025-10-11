package com.example.point.batch.config;

import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.jdbc.DataSourceBuilder;
import org.springframework.boot.orm.jpa.EntityManagerFactoryBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;
import org.springframework.orm.jpa.JpaTransactionManager;
import org.springframework.orm.jpa.LocalContainerEntityManagerFactoryBean;
import org.springframework.transaction.PlatformTransactionManager;

import javax.sql.DataSource;

@Configuration
@EnableJpaRepositories(
        basePackages = "com.example.point.batch.point.repos",
        entityManagerFactoryRef = "pointEntityManagerFactory",
        transactionManagerRef = "pointTransactionManager"
)
public class PointDataSourceConfig {

    @Bean
    @ConfigurationProperties("spring.datasource.point")
    public DataSource pointDataSource() {
        return DataSourceBuilder.create().build();
    }

    @Bean(name = "pointEntityManagerFactory")
    public LocalContainerEntityManagerFactoryBean pointEntityManagerFactory(
            EntityManagerFactoryBuilder builder,
            @Qualifier("pointDataSource") DataSource dataSource
    ) {
        return builder
                .dataSource(dataSource)
                .packages("com.example.point.batch.point.domain")
                .persistenceUnit("pointPU")
                .build();
    }

    @Bean(name = "pointTransactionManager")
    public PlatformTransactionManager pointTransactionManager(
            @Qualifier("pointEntityManagerFactory") LocalContainerEntityManagerFactoryBean emf
    ) {
        return new JpaTransactionManager(emf.getObject());
    }
}

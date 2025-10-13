package com.example.catalog.cassandra;

import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.cassandra.core.CassandraAdminTemplate;
import org.springframework.data.cassandra.core.cql.CqlTemplate;

@Configuration
public class CassandraConfig {

    private static final String KEYSPACE_NAME = "catalog";

    @Bean
    @ConditionalOnMissingBean
    public Boolean createKeyspaceIfNotExists(CqlTemplate cqlTemplate) {
        String createKeyspace =
            "CREATE KEYSPACE IF NOT EXISTS " + KEYSPACE_NAME + " " +
            "WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}";
        cqlTemplate.execute(createKeyspace);
        return true; // just to ensure bean creation
    }
}

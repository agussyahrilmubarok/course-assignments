package com.example.catalog.cassandra;

import org.springframework.boot.autoconfigure.domain.EntityScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.cassandra.repository.config.EnableCassandraRepositories;

@Configuration
@EnableCassandraRepositories("com.example.catalog.cassandra.repos")
@EntityScan("com.example.catalog.cassandra.domain")
public class CassandraConfig {

}

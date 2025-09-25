package com.finly.finly.repos.impl;

import com.finly.finly.domain.Customer;
import com.finly.finly.repos.CustomerCustomRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.data.mongodb.core.query.Criteria;
import org.springframework.data.mongodb.core.query.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
@RequiredArgsConstructor
public class CustomerCustomRepositoryImpl implements CustomerCustomRepository {

    private final MongoTemplate mongoTemplate;

    @Override
    public List<Customer> findCustomersByUserAndSearch(UUID userId, String search) {
        Query query = new Query();
        //query.addCriteria(Criteria.where("user.$id").is(userId));
        query.addCriteria(Criteria.where("user").is(userId));
        if (search != null && !search.isEmpty()) {
            Criteria searchCriteria = new Criteria().orOperator(
                    Criteria.where("name").regex(search, "i"),
                    Criteria.where("email").regex(search, "i"),
                    Criteria.where("phone").regex(search, "i"),
                    Criteria.where("address").regex(search, "i")
            );
            query.addCriteria(searchCriteria);
        }

        return mongoTemplate.find(query, Customer.class);
    }

    @Override
    public List<Customer> findCustomersByUserId(UUID userId) {
        Query query = new Query();
        //query.addCriteria(Criteria.where("user.$id").is(userId));
        query.addCriteria(Criteria.where("user").is(userId));
        return mongoTemplate.find(query, Customer.class);
    }
}

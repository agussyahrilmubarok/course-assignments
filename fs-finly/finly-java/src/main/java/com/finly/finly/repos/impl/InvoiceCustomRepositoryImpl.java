package com.finly.finly.repos.impl;

import com.finly.finly.domain.Invoice;
import com.finly.finly.repos.InvoiceCustomRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.data.mongodb.core.query.Criteria;
import org.springframework.data.mongodb.core.query.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
@RequiredArgsConstructor
public class InvoiceCustomRepositoryImpl implements InvoiceCustomRepository {

    private final MongoTemplate mongoTemplate;

    @Override
    public List<Invoice> findInvoicesByOwnerAndSearch(UUID ownerId, String search) {
        Query query = new Query();
        //query.addCriteria(Criteria.where("owner.$id").is(ownerId));
        query.addCriteria(Criteria.where("owner").is(ownerId));
        if (search != null && !search.isEmpty()) {
            Criteria searchCriteria = new Criteria().orOperator(
                    Criteria.where("status").regex(search, "i"),
                    Criteria.where("customer.name").regex(search, "i"),
                    Criteria.where("customer.email").regex(search, "i"),
                    Criteria.where("customer.phone").regex(search, "i")
            );
            query.addCriteria(searchCriteria);
        }

        return mongoTemplate.find(query, Invoice.class);
    }

    @Override
    public List<Invoice> findInvoicesByOwnerId(UUID ownerId) {
        Query query = new Query();
        //query.addCriteria(Criteria.where("owner.$id").is(ownerId));
        query.addCriteria(Criteria.where("owner").is(ownerId));
        return mongoTemplate.find(query, Invoice.class);
    }
}

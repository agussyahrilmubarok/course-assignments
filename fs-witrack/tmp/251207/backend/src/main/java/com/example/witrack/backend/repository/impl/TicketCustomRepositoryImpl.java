package com.example.witrack.backend.repository.impl;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.repository.TicketCustomRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.data.mongodb.core.query.Criteria;
import org.springframework.data.mongodb.core.query.Query;
import org.springframework.stereotype.Repository;

import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.List;

@Repository
@RequiredArgsConstructor
public class TicketCustomRepositoryImpl implements TicketCustomRepository {

    private final MongoTemplate mongoTemplate;

    @Override
    public List<Ticket> searchTickets(String keyword,
                                      Ticket.Status status,
                                      Ticket.Priority priority,
                                      OffsetDateTime startAt,
                                      OffsetDateTime endAt) {
        List<Criteria> criteriaList = new ArrayList<>();

        if (keyword != null && !keyword.isBlank()) {
            criteriaList.add(new Criteria().orOperator(
                    Criteria.where("title").regex(keyword, "i"),
                    Criteria.where("description").regex(keyword, "i")
            ));
        }

        if (status != null) {
            criteriaList.add(Criteria.where("status").is(status));
        }

        if (priority != null) {
            criteriaList.add(Criteria.where("priority").is(priority));
        }

        if (startAt != null && endAt != null) {
            criteriaList.add(Criteria.where("createdAt").gte(startAt).lte(endAt));
        } else if (startAt != null) {
            criteriaList.add(Criteria.where("createdAt").gte(startAt));
        } else if (endAt != null) {
            criteriaList.add(Criteria.where("createdAt").lte(endAt));
        }

        Query query = new Query();

        if (!criteriaList.isEmpty()) {
            query.addCriteria(new Criteria().andOperator(criteriaList.toArray(new Criteria[0])));
        }

        return mongoTemplate.find(query, Ticket.class);
    }

    @Override
    public List<Ticket> searchTickets(String keyword,
                                      Ticket.Status status,
                                      Ticket.Priority priority,
                                      OffsetDateTime startAt,
                                      OffsetDateTime endAt,
                                      String userID) {
        List<Criteria> criteriaList = new ArrayList<>();

        if (keyword != null && !keyword.isBlank()) {
            criteriaList.add(new Criteria().orOperator(
                    Criteria.where("title").regex(keyword, "i"),
                    Criteria.where("description").regex(keyword, "i")
            ));
        }

        if (status != null) {
            criteriaList.add(Criteria.where("status").is(status));
        }

        if (priority != null) {
            criteriaList.add(Criteria.where("priority").is(priority));
        }

        if (startAt != null && endAt != null) {
            criteriaList.add(Criteria.where("createdAt").gte(startAt).lte(endAt));
        } else if (startAt != null) {
            criteriaList.add(Criteria.where("createdAt").gte(startAt));
        } else if (endAt != null) {
            criteriaList.add(Criteria.where("createdAt").lte(endAt));
        }

        if (userID != null && !userID.isBlank()) {
            criteriaList.add(Criteria.where("user").is(userID));
        }

        Query query = new Query();

        if (!criteriaList.isEmpty()) {
            query.addCriteria(new Criteria().andOperator(criteriaList.toArray(new Criteria[0])));
        }

        return mongoTemplate.find(query, Ticket.class);
    }
}

package com.example.grpc.client.service;

import com.example.bookstore.*;
import com.example.grpc.client.model.ReviewDTO;
import com.google.protobuf.Empty;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.client.inject.GrpcClient;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

@Service
@Slf4j
public class ReviewService {

    private final ReviewServiceGrpc.ReviewServiceBlockingStub stub;

    public ReviewService(@GrpcClient("book-server") ReviewServiceGrpc.ReviewServiceBlockingStub stub) {
        this.stub = stub;
    }

    public List<ReviewDTO.Response> findAll() {
        log.info("Fetching all reviews from gRPC server...");
        try {
            Iterator<Review> reviewIterator = stub.findAll(Empty.newBuilder().build());
            List<ReviewDTO.Response> reviews = new ArrayList<>();
            while (reviewIterator.hasNext()) {
                reviews.add(ReviewDTO.Response.from(reviewIterator.next()));
            }
            log.info("Fetched {} reviews successfully", reviews.size());
            return reviews;
        } catch (Exception e) {
            log.error("Failed to fetch reviews: {}", e.getMessage(), e);
            throw e;
        }
    }

    public ReviewDTO.Response get(Long id) {
        log.info("Fetching review with id={} from gRPC server...", id);
        try {
            Review protoReview = stub.get(GetReviewRequest.newBuilder().setId(id).build());
            log.info("Fetched review: id={}, content='{}'", protoReview.getId(), protoReview.getContent());
            return ReviewDTO.Response.from(protoReview);
        } catch (Exception e) {
            log.error("Failed to fetch review with id={}: {}", id, e.getMessage(), e);
            throw e;
        }
    }

    public ReviewDTO.Response create(ReviewDTO.Request request) {
        log.info("Creating new review for bookId={} with rating={}...", request.getBookId(), request.getRating());
        try {
            CreateReviewRequest protoRequest = CreateReviewRequest.newBuilder()
                    .setContent(request.getContent())
                    .setRating(request.getRating().floatValue())
                    .setBookId(request.getBookId())
                    .build();
            Review protoReview = stub.create(protoRequest);
            log.info("Review created: id={}, content='{}'", protoReview.getId(), protoReview.getContent());
            return ReviewDTO.Response.from(protoReview);
        } catch (Exception e) {
            log.error("Failed to create review for bookId={}: {}", request.getBookId(), e.getMessage(), e);
            throw e;
        }
    }

    public ReviewDTO.Response update(Long id, ReviewDTO.Request request) {
        log.info("Updating review id={} with new content='{}' and rating={}...", id, request.getContent(), request.getRating());
        try {
            UpdateReviewRequest protoRequest = UpdateReviewRequest.newBuilder()
                    .setId(id)
                    .setContent(request.getContent())
                    .setRating(request.getRating().floatValue())
                    .setBookId(request.getBookId())
                    .build();
            Review protoReview = stub.update(protoRequest);
            log.info("Review updated: id={}, content='{}'", protoReview.getId(), protoReview.getContent());
            return ReviewDTO.Response.from(protoReview);
        } catch (Exception e) {
            log.error("Failed to update review id={}: {}", id, e.getMessage(), e);
            throw e;
        }
    }

    public void delete(Long id) {
        log.info("Deleting review with id={}...", id);
        try {
            stub.delete(GetReviewRequest.newBuilder().setId(id).build());
            log.info("Review with id={} deleted successfully", id);
        } catch (Exception e) {
            log.error("Failed to delete review id={}: {}", id, e.getMessage(), e);
            throw e;
        }
    }
}
